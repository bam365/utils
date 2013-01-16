/* status.go - Continuously outputs status for status monitor
 * Blake Mitchell, 2011
 */

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Structs

type Atom struct {
	Ticks     int
	Label     string
	Value     func() string
	TicksLeft int
	Output    string
}

type ColorMarkup struct {
	Lcbegin string
	Lcend   string
	Vcbegin string
	Vcend   string
}

// Methods

func (a *Atom) set_labeled_output(cm *ColorMarkup) {
	a.Output = ""
	a.Output = cm.Lcbegin + "[" + a.Label + ":" + cm.Lcend
	a.Output += cm.Vcbegin + a.Value() + cm.Vcend + cm.Lcbegin + "]"
	a.Output += cm.Lcend
}

// Main

func main() {
	atoms := []Atom{{1, "C", get_cpu_usage(), 0, ""},
		{5, "M", get_mem_usage, 0, ""},
		{2, "N", get_net_usage("eth0", 2), 0, ""},
		{2, "VOL", get_volume("Master"), 0, ""},
		{5 * 60, "GM", get_gmail, 0, ""},
		{60, "PA", get_pkg_update_num, 0, ""},
		{30, "", get_time, 0, ""},
	}
	color := ColorMarkup{os.Getenv("STATCL1"), os.Getenv("STATCLEND"),
		os.Getenv("STATCL2"), os.Getenv("STATCLEND")}
	var out string
	prefix := os.Getenv("STATPREFIX")
	suffix := os.Getenv("STATSUFFIX")
	for {
		out = generate_status(atoms, " ", &color)
		out = prefix + out + suffix
		fmt.Println(out)
		time.Sleep(1e9)
	}
}

// Functions

func generate_status(atoms []Atom, sep string, cm *ColorMarkup) string {
	out := ""
	for i := 0; i < len(atoms); i++ {
		if atoms[i].TicksLeft <= 0 {
			if atoms[i].Label != "" {
				atoms[i].set_labeled_output(cm)
			} else {
				atoms[i].Output = cm.Vcbegin + atoms[i].Value()
				atoms[i].Output += cm.Vcend
			}
			atoms[i].TicksLeft = atoms[i].Ticks
		}
		atoms[i].TicksLeft--
		out += atoms[i].Output
		if i < len(atoms)-1 {
			out += sep
		}
	}
	return out
}

func str_pad(s string, p string, n int) string {
	ret := s
	if len(s) < n {
		ret = strings.Repeat(p[:1], n-len(s)) + s
	}
	return ret
}

func get_output(cmd *exec.Cmd) string {
	buf, err := cmd.CombinedOutput()
	if err != nil {
		return "?"
	}
	out := string(buf)
	return out
}

func get_field(s string, n int) string {
	f := strings.Fields(s)
	if n >= len(f) {
		return ""
	}
	return f[n]
}

func xfer_rate_str(rate int) string {
	var out string
	var suf string

	if rate >= 1000000000 {
		out = strconv.Itoa(rate / 1000000000)
		suf = "G"
	} else if rate >= 10000000 {
		out = strconv.Itoa(rate / 1000000)
		suf = "M"
	} else if rate >= 1000000 {
		out = strconv.FormatFloat(float64(float32(rate)/1000000.0), 'f', 2, 32)
		suf = "M"
	} else if rate >= 1000 {
		out = strconv.Itoa(rate / 1000)
		suf = "K"
	} else {
		out = strconv.Itoa(rate)
		suf = "B"
	}
	if len(out) > 3 {
		out = out[:3]
	} else if len(out) < 3 {
		out = str_pad(out, " ", 3)
	}
	return out + suf
}

func get_cpu_usage() func() string {
	lastidle := 0
	lasttotal := 0
	return func() string {
		var ret string

		out := get_output(exec.Command("cat", "/proc/stat"))
		newtotal := 0
		for i := 1; i < 5; i++ {
			a, _ := strconv.Atoi(get_field(out, i))
			newtotal += a
		}
		newidle, _ := strconv.Atoi(get_field(out, 4))
		diffidle := newidle - lastidle
		difftotal := newtotal - lasttotal
		lastidle = newidle
		lasttotal = newtotal
		if perc := 99 - (diffidle*100)/difftotal; perc < 0 {
			ret = "0"
		} else {
			ret = strconv.Itoa(perc)
		}
		return str_pad(ret, " ", 2)
	}
}

func get_mem_usage() string {
	out := get_output(exec.Command("free"))
	used, _ := strconv.Atoi(get_field(out, 15))
	total, _ := strconv.Atoi(get_field(out, 7))
	ret := used * 100 / total
	return str_pad(strconv.Itoa(ret), " ", 2)
}

func get_net_usage(intf string, loadavg int) func() string {
	var lastup, lastdown uint64 = 0, 0
	return func() string {
		nstat := get_output(exec.Command("cat", "/proc/net/dev"))
		nstat = nstat[strings.Index(nstat, intf)+len(intf)+1:]
		newup, _ := strconv.ParseUint(get_field(nstat, 8), 10, 64)
		newdown, _ := strconv.ParseUint(get_field(nstat, 0), 10, 64)
		uprate := int((newup - lastup) / 2)
		downrate := int((newdown - lastdown) / 2)
		out := xfer_rate_str(uprate) + "/" + xfer_rate_str(downrate)
		lastup, lastdown = newup, newdown
		return out
	}
}

func get_volume(mixer string) func() string {
	return func() string {
		out := get_output(exec.Command("volleft.sh", mixer))
		return str_pad(strings.Trim(out, " \n"), " ", 3)
	}
}

func get_gmail() string {
	out := get_output(exec.Command("checkgmail.py"))
	return str_pad(strings.Trim(out, " \n"), " ", 2)
}

func get_pkg_update_num() string {
	out := get_output(exec.Command("checkpacman.py"))
	return str_pad(strings.Trim(out, " \n"), " ", 3)
}

func get_time() string {
	t := time.Now()
	var out, ampm, hr string

	if t.Hour() >= 12 {
		ampm = "PM"
	} else {
		ampm = "AM"
	}
	hr = strconv.Itoa(t.Hour() % 12)
	if hr == "0" {
		hr = "12"
	}
	min := str_pad(strconv.Itoa(t.Minute()), "0", 2)
	out = str_pad(hr, " ", 2) + ":" + min + " " + ampm

	return out
}
