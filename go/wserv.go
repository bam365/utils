/* wserv.go - Single-page webserver that serves the given page on localhost/
* Blake Mitchell, 2013
*/

package main

import (
        "fmt"
        "flag"
        "io/ioutil"
        "net/http"
)


func main() {
        var portflag int
        var fcontents string

        flag.IntVar(&portflag, "p", 8080, "Port to run webserver on")
        flag.Parse()
        if len(flag.Args()) < 1 {
                fmt.Println("Need to specify a page")
                return;
        }
        fname := flag.Arg(0);

        fbuf, ferr := ioutil.ReadFile(fname)
        if ferr != nil {
                fmt.Printf("Could not open file. Reason: %s", ferr.Error())
                return
        }
        fcontents = string(fbuf)

        fmt.Printf("Serving %s on port %d\n", fname, portflag)
        http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
                fmt.Fprintln(w, fcontents)
        });
        herr := http.ListenAndServe(fmt.Sprintf(":%d", portflag), nil)
        if herr != nil {
                fmt.Printf("Couldn't start web server. Reason: %s\n", 
                            herr.Error())
        }

}






               



