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
var shell string
var command string

func init() {
	flag.StringVar(&shell, "shell", "bash", "Shell to run command with")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal(Usage)
	}
	command = args[0]
	// prime the file list with zero'd times
	for _, v := range args[1:] {
		filelist[v] = time.Time{}
	}
}

func timestamp() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func updateFileModTimes() time.Time {
	var maxTime time.Time
	if len(filelist) == 0 {
		log.Fatal("No files to monitor.")
	}
	for k, _ := range filelist {
		s, err := os.Stat(k)
		for i := 3; i > 0 && err != nil; i-- {
			time.Sleep(100 * time.Millisecond)
			s, err = os.Stat(k)
		}
		if err != nil {
			log.Println("Removing file: ", k)
			delete(filelist, k)
			maxTime = time.Now()
		} else {
			filelist[k] = s.ModTime()
			if filelist[k].After(maxTime) {
				maxTime = filelist[k]
			}
		}
	}
	return maxTime
}

func main() {
	var lastTime time.Time
	var currTime time.Time

	for {
		currTime = updateFileModTimes()
		if lastTime.Before(currTime) {
			fmt.Println(timestamp(), "Begin")
			cmd := exec.Command(shell, "-c", command)
			// ignore error as we want it reported and keep running
			out, _ := cmd.CombinedOutput()
			fmt.Println(strings.TrimRight(string(out), "\n "))
			fmt.Println(timestamp(), "End - Monitoring", len(filelist), "files.")
		}
		// pause to not overrun the system
		time.Sleep(500 * time.Millisecond)
		lastTime = currTime
	}
}
