/* wserv.go - Single-page webserver that serves the given page on localhost/
* Blake Mitchell, 2013
*/

package main

import (
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
        "strconv"
)


func main() {
        port := 8080

        if len(os.Args) > 1 {
                if pparse, cerr := strconv.Atoi(os.Args[1]); cerr != nil {
                        fmt.Println("Bad port arg")
                        return;
                } else {
                        port = pparse
                }
        }

        fmt.Printf("Serving on port %d\n", port)
        herr := http.ListenAndServe(fmt.Sprintf(":%d", port),
                                    http.FileServer(http.Dir(".")))
        if herr != nil {
                fmt.Printf("Couldn't start web server. Reason: %s\n", 
                            herr.Error())
        }

}


func Handler(w http.ResponseWriter, r *http.Request) {
        fname := r.URL.RequestURI()

        if fname[0] == '/' {
                fname = fname[1:]
        }
        fbuf, ferr := ioutil.ReadFile(fname)
        if ferr != nil {
                fmt.Printf("Could not open file '%s'. Reason: %s", fname, ferr.Error())
                http.Error(w, "Not found", http.StatusNotFound)
                return
        } else {
                fmt.Printf("File '%s' requested\n", fname)
        }

        fmt.Fprintln(w, string(fbuf))
}

        
        
