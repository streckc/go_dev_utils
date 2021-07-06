package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"time"
)

var Version string
var Build string

var inactivity int
var shell string
var version bool
var remove bool
var filelist = make(map[string]time.Time)
var command string

var lastRun time.Time

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] <command> <file> [files..]\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "\nPositional:\n")
	fmt.Fprintf(os.Stderr, "  command - String of command to run when files change\n")
	fmt.Fprintf(os.Stderr, "  file    - File to monitor\n")
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
}

func displayVersion() {
	fmt.Printf("Version: %v\n", Version)
	fmt.Printf("Build:   %v\n", Build)
}

func setupArgs() {
	flag.Usage = Usage
	flag.BoolVar(&version, "V", false, "Version: Display version and build")
	flag.BoolVar(&remove, "R", false, "Remove: Remove files from monitor list if not found")
	flag.IntVar(&inactivity, "i", 600, "Inactivity: Seconds of inactivity before exiting")
	flag.StringVar(&shell, "s", "bash", "Shell: Shell to run command with")
	flag.Parse()
}

func init() {
	setupArgs()

	if version {
		displayVersion()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		log.Fatal("missing required parameters.")
	}
	command = args[0]

	// prime the file list with zero'd times
	for _, v := range args[1:] {
		filelist[v] = time.Time{}
	}
}

// updatefileModTimes will stat each file being monitored and updtae the list
// with the files modified time.
func updateFileModTimes() time.Time {
	var maxTime time.Time
	// Track missing files to no report time if files are missing.
	missing := 0

	// If no files are being monitored, then exit.
	if len(filelist) == 0 {
		log.Fatal("No files to monitor.")
	}

	// Range over the file list to get the updated times.
	for k, t := range filelist {
		s, err := os.Stat(k)
		// Account for files disappearing temporarily during writes
		for i := 3; i > 0 && err != nil; i-- {
			time.Sleep(100 * time.Millisecond)
			s, err = os.Stat(k)
		}
		if err != nil {
			if remove {
				log.Println("Removing file: ", k)
				delete(filelist, k)
				maxTime = time.Now()
			}
			if t.After(time.Time{}) {
				log.Println("Missing file: ", k)
			}
			filelist[k] = time.Time{}
			missing++
		} else {
			newTime := s.ModTime()
			if t.Equal(time.Time{}) {
				log.Println("Adding file: ", k)
			} else if newTime.After(t) {
				log.Println("Updating file: ", k)
			}
			if newTime.After(maxTime) {
				maxTime = newTime
			}
			filelist[k] = newTime
		}
	}

	// if missing files, zero out maxTime
	if missing > 0 {
		maxTime = time.Time{}
	}

	return maxTime
}

func monitorInactivity() {
	// Inactivity of 0 will run until interrupted
	if inactivity <= 0 {
		return
	}
	for {
		currTime := time.Now().Add(time.Duration(-1*inactivity) * time.Second)
		// Only check inactivity after first run of the command completes
		if !lastRun.IsZero() && currTime.After(lastRun) {
			log.Fatal("Inactive for ", inactivity, " seconds.  Exiting.")
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	var currTime time.Time
	var lastTime time.Time

	go monitorInactivity()

	for {
		currTime = updateFileModTimes()
		if lastTime.Before(currTime) {
			log.Println("Begin - Running command.")

			var stdBuffer bytes.Buffer
			cmd := exec.Command(shell, "-c", command)
			mw := io.MultiWriter(os.Stdout, &stdBuffer)
			cmd.Stdout = mw
			cmd.Stderr = mw

			cmd.Run()
			//log.Println(stdBuffer.String())
			log.Println("End - Monitoring", len(filelist), "files.")
			lastRun = time.Now()
		}
		// pause to not overrun the system
		time.Sleep(500 * time.Millisecond)
		lastTime = currTime
	}
}
