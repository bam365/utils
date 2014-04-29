open Core.Std


module Tag: sig
    type t

    val from_str: string -> t option
    val begin_str: t -> string
    val end_str: t -> string

end = struct
    type t = {
        tag_name: string;
        id: string option;
        clss: string option;
    }

    let tag_str_is_valid s =
        let char_count c = String.count s ~f:((=) c) in
        (char_count '.') <= 1 && (char_count '#') <= 1 

    let get_tag_name s = 
        match String.split_on_chars s ~on:['#'; '.'] with
        | tag_name::_ -> tag_name
        | _ -> "" (* should never happen *)

    let get_attribute s first_char last_char =
        match String.split s ~on:first_char with
        | _::first_str::_ ->
                (match String.split first_str ~on:last_char with
                | str::_ -> Some str
                | _ -> Some first_str)
        | _ -> None


    let from_str s =
        if tag_str_is_valid s then
            Some { tag_name = get_tag_name s; 
                   id = get_attribute s '#' '.'; 
                   clss = get_attribute s '.' '#';
                 }
        else None

    let begin_str t =
        let attr_str attr_name attr  = 
            match attr with
            | Some s -> sprintf " %s=\"%s\"" attr_name s
            | None -> ""
        in
        sprintf "<%s%s%s>" t.tag_name (attr_str "id" t.id) 
                           (attr_str "class" t.clss)

    let end_str t = sprintf "</%s>" t.tag_name
end
    


let trim_split str =
    let not_whitespace _ c = Char.is_whitespace c |> not in
    let str_len = String.length str in
    let left_pos = 
        String.lfindi str ~f:not_whitespace |> Option.value ~default:0
    in
    let right_pos =
        String.rfindi str ~f:not_whitespace 
        |> Option.value ~default:(str_len - 1)
        |> (+) 1
    in
    let part p q = String.slice str p q in
    let left_str, right_str = 
        (if left_pos = 0 then "" else part 0 left_pos),
        (if right_pos = str_len then "" else String.slice str right_pos 0)
        
    in
    let mid_str = part left_pos right_pos in
    [left_str; mid_str; right_str]


let add_tag_between_ws is_block tag str =
    let [left_ws; mid; right_ws] = trim_split str in
    let tag_sep = if is_block then "\n" ^ left_ws else "" in
    sprintf "%s%s%s%s%s%s%s" 
        left_ws 
        (Tag.begin_str tag) tag_sep mid tag_sep (Tag.end_str tag)
        right_ws

let main is_inline tag_str () =
    let tag_str = Option.value tag_str ~default:"" in 
    match Tag.from_str tag_str with
    | Some tag ->
        let transform = add_tag_between_ws in
        In_channel.input_all In_channel.stdin
        |> transform (not is_inline) tag
        |> Out_channel.output_string Out_channel.stdout
    | None -> 
        "Malformed tag"
        |> Out_channel.output_string Out_channel.stderr;
        Out_channel.output_string Out_channel.stdout tag_str


(* Command stuff *)

let spec =
    let open Command.Spec in
    empty 
    +> flag "-i" no_arg ~doc:"inline tag (no newlines after tags)"
    +> anon (maybe ("tag" %: string))

let command =
    Command.basic
        ~summary:"Puts given text into given tag"
        ~readme:(fun () -> "None yet")
        spec
        main

let () = Command.run ~version:"1.0" ~build_info:"bam" command
