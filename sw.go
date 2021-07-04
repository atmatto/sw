// Minimal stopwatch
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func usage() {
	fmt.Println("usage: sw reset/resume/pause/show")
}

func displayTime(u int64) string {
	str := ""
	if len(strconv.FormatInt((u/60), 10)) == 1 {
		str = string(append([]byte(str), '0'))
	}
	str = string(append([]byte(str), []byte(strconv.FormatInt((u/60), 10))...))
	str = string(append([]byte(str), ':'))
	if len(strconv.FormatInt(u%60, 10)) == 1 {
		str = string(append([]byte(str), '0'))
	}
	str = string(append([]byte(str), []byte(strconv.FormatInt(u%60, 10))...))
	return str
}

func main() {
	a := os.Args
	if len(a) == 1 {
		usage()
		return
	}
	mode := a[1]

	h, err := os.UserHomeDir()
	if err != nil {
		log.Panicln("sw: couldn't get user home directory:", err)
	}
	sw, err := os.OpenFile(filepath.Join(h, "/.sw"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Panicln("sw: couldn't open ~/.sw file:", err)
	}
	defer sw.Close()

	swc, err := io.ReadAll(sw)
	if err != nil {
		log.Panicln("sw: couldn't read ~/.sw file:", err)
	}

	switch mode {
	case "reset":
		err = sw.Truncate(0)
		if err != nil {
			log.Panicln("sw: couldn't truncate ~/.sw file:", err)
		}
		sw.Seek(0, 0)
		fmt.Fprint(sw, "z")
	case "resume":
		if swc[0] == 'z' {
			err = sw.Truncate(0)
			if err != nil {
				log.Panicln("sw: couldn't truncate ~/.sw file:", err)
			}
			sw.Seek(0, 0)
			fmt.Fprint(sw, time.Now().Unix())
		} else if swc[0] == 'p' {
			et, err := strconv.ParseInt(string(swc[1:]), 10, 64)
			if err != nil {
				log.Panicln("sw: couldn't parse ~/.sw file:", err)
			}
			st := time.Now().Unix() - et
			err = sw.Truncate(0)
			if err != nil {
				log.Panicln("sw: couldn't truncate ~/.sw file:", err)
			}
			sw.Seek(0, 0)
			fmt.Fprintf(sw, "%d", st)
		}
	case "pause":
		if swc[0] != 'z' && swc[0] != 'p' {
			t, err := strconv.ParseInt(string(swc), 10, 64)
			if err != nil {
				log.Panicln("sw: couldn't parse ~/.sw file:", err)
			}
			et := time.Now().Unix() - t
			err = sw.Truncate(0)
			if err != nil {
				log.Panicln("sw: couldn't truncate ~/.sw file:", err)
			}
			sw.Seek(0, 0)
			fmt.Fprintf(sw, "p%d", et)
		}
	case "show":

	default:
		usage()
	}
	if swc[0] == 'p' {
		t, err := strconv.ParseInt(string(swc[1:]), 10, 64)
		if err != nil {
			log.Panicln("sw: couldn't parse ~/.sw file:", err)
		}
		fmt.Println(displayTime(t))
	} else if swc[0] == 'z' {
		fmt.Println("--:--")
	} else {
		t, err := strconv.ParseInt(string(swc), 10, 64)
		if err != nil {
			log.Panicln("sw: couldn't parse ~/.sw file:", err)
		}
		fmt.Println(displayTime(time.Now().Unix() - t))
	}
}
