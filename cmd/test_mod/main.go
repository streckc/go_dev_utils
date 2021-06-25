package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const Usage = "test_mod <command> <file> [<file>...]"

var filelist = make(map[string]time.Time)
var command []string

func init() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal(Usage)
	}
	command = strings.Split(args[0], " ")

	// prime the file list with zero'd times
	for _, v := range args[1:] {
		filelist[v] = time.Time{}
	}
}

func timestamp() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func getMaxFileTime() time.Time {
	var maxTime time.Time

	for _, t := range filelist {
		if t.After(maxTime) {
			maxTime = t
		}
	}

	return maxTime
}

func updateFileModTimes() {
	for k, _ := range filelist {
		s, err := os.Stat(k)
		for i := 3; i > 0 && err != nil; i-- {
			time.Sleep(100 * time.Millisecond)
			s, err = os.Stat(k)
		}
		if err != nil {
			log.Fatal(err)
		}
		filelist[k] = s.ModTime()
	}
}

func main() {
	var lastTime time.Time
	var currTime time.Time

	for {
		updateFileModTimes()
		currTime = getMaxFileTime()
		if lastTime.Before(currTime) {
			fmt.Println(timestamp(), "Begin")
			out, err := exec.Command(command[0], command[1:]...).Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(strings.TrimRight(string(out), "\n "))
			fmt.Println(timestamp(), "End")
		}
		// pause to not overrun the system
		time.Sleep(500 * time.Millisecond)
		lastTime = currTime
	}
}
