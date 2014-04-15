open Core.Std


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
    sprintf "%s<%s>%s%s%s</%s>%s" 
        left_ws 
        tag tag_sep mid tag_sep tag
        right_ws

let main is_block tag () =
    let transform = add_tag_between_ws in
    In_channel.input_all In_channel.stdin
    |> transform is_block (Option.value tag ~default:"")
    |> Out_channel.output_string Out_channel.stdout


(* Command stuff *)

let spec =
    let open Command.Spec in
    empty 
    +> flag "-b" no_arg ~doc:"block tag (put newlines after tags)"
    +> anon (maybe ("tag" %: string))

let command =
    Command.basic
        ~summary:"Puts given text into given tag"
        ~readme:(fun () -> "None yet")
        spec
        main

let () = Command.run ~version:"1.0" ~build_info:"bam" command
