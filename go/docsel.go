// docsel.go - Use dmenu to select a document
// Blake Mitchell, 2012 

package main

import (
        "fmt"
        "os"
        "bufio"
        "strings"
        "exec"
        "unicode"
        "flag"
)

type DocselArgs struct {
        fn, nb, nf, sb, sf  string
}

func (da *DocselArgs) Parse() {
        flag.StringVar(&da.fn, "fn", "", "Font")
        flag.StringVar(&da.nb, "nb", "black", "Normal background color")
        flag.StringVar(&da.nf, "nf", "green", "Normal foreground color")
        flag.StringVar(&da.sb, "sb", "green", "Selected background color")
        flag.StringVar(&da.sf, "sf", "black", "Selected foreground color")
        flag.Parse()
}


func main() {
        da := new(DocselArgs)
        da.Parse()
        lines := readlines(bufio.NewReader(os.Stdin));
        doc := selectdoc(lines, da)
        editdoc(doc)
}


func readlines(s *bufio.Reader) []string {
        lines := make([]string, 0, 1)
        for  buf, err := s.ReadString('\n'); err == nil; 
             buf, err = s.ReadString('\n') {
                lines = append(lines, buf)
        }

        return lines
}

              
func getfilenames(files []string) []string {
        pn := make([]string, 0);
        for _, s := range(files) {
                pn = append(pn, getname(s))
        }
                        
        return pn
}


func getfilefromname(files []string, name string) string {
        var fulldoc string
        for _, s := range(files) {
                if getname(s) == name {
                        fulldoc = s
                }
        }
        return fulldoc
}


func getname(line string) string {
        var ret string
        if i := strings.LastIndex(line, "/"); i != -1 && i < len(line)-1 {
                ret = line[i+1:]
        } else {
                ret = line
        }
        return strings.TrimRight(ret, "\n ")
}


func getdir(line string) string {
        var ret string
        if i := strings.LastIndex(line, "/"); i != -1 {
                ret =  line[:i]
        } else {
                ret = ""
        }
        return strings.TrimSpace(ret)
}


func selectdoc(files []string, da *DocselArgs) string {
        dcmd := exec.Command("dmenu", "-nb", da.nb, "-nf", da.nf,
                             "-sf", da.sf, "-sb", da.sb, "-fn", da.fn)
        din, _ := dcmd.StdinPipe()
        dout, _ := dcmd.StdoutPipe()
        rbuf := make([]byte, 512)
        dcmd.Start();
        for _, doc := range(getfilenames(files)) {
                din.Write([]byte(doc))
                din.Write([]byte{'\n'})
        }
        din.Close()
        dout.Read(rbuf)
        dcmd.Wait()
        doc := strings.TrimFunc(string(rbuf), 
                func(r int) bool {return unicode.IsControl(r) || 
                                  unicode.IsSpace(r)})
        return getfilefromname(files, doc)
}


func editdoc(line string) {
        dir, name := getdir(line), getname(line)
        if name == "" {
                return;
        }
        os.Chdir(os.Getenv("HOME"))
        if dir != "" { 
                if err := os.Chdir(dir); err != nil {
                        fmt.Println("error: " + err.String())
                }
        }
        vimpath, _ := exec.LookPath("gvim")
        args := []string{vimpath, name}
        if err := os.Exec(vimpath, args, os.Environ()); err != nil {
                fmt.Println("error: " + err.String())
        }
}

