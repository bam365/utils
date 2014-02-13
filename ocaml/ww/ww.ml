open Core.Std

(* Ocaml's sucky regexs don't have a whitespace match pattern *)
let words = 
    let whitespace_regexp =
        let whitespace_chars =
            [9; 10; 11; 12; 13; 32]
            |> List.map ~f:(fun c -> c |> Char.of_int_exn |> String.make 1)
            |> String.concat 
        in
        Str.regexp ("[" ^ whitespace_chars ^ "]+")
    in
    Str.split whitespace_regexp

let unwords words' = 
    List.map ~f:(String.concat ~sep:" ") words' |> String.concat ~sep:"\n"

let append_word_to_last_line words' w =
    match words' with
    | line::rest -> (w::line)::rest
    | [] -> [[w]]

let wrap_words width wrds  =
    let rec loop length acc words' =
        match words' with
        (* we've built this up in reverse order *)
        | [] -> List.rev (List.map ~f:List.rev acc) 
        | w::rest ->
            let word_length = String.length w in
            let (new_length, new_acc) =
                if (length + word_length) > width then
                    (word_length + 1), ([w]::acc)
                else
                    (length + word_length + 1),
                    (append_word_to_last_line acc w) 
            in
            loop new_length new_acc rest
    in
    loop 0 [] wrds


let wrap_text width text = text |> words |> wrap_words width |> unwords

let wrap_channel width in_chan out_chan =
    In_channel.input_all in_chan
    |> wrap_text width
    |> Out_channel.output_string out_chan;
    Out_channel.newline out_chan


(* Command stuff *)

let do_wrap width filename () =
    let in_chan = 
        match filename with
        | "-" -> In_channel.stdin
        | filename -> In_channel.create filename
    in
    wrap_channel width in_chan Out_channel.stdout

let ww_command =
    Command.basic
        ~summary:"Wraps text on stdin to stdout"
        Command.Spec.(
            empty 
            +> flag "-w" (optional_with_default 80 int) ~doc:"Width of wrapping"
            +> anon (maybe_with_default "-" ("filename" %: file)))
        do_wrap

let () = Command.run ww_command
