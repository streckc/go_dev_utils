package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"time"
)

var Version string
var Build string
var Date string

var version bool
var commands bool
var location string
var archive string
var days int

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] <location>\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "\nPositional:\n")
	fmt.Fprintf(os.Stderr, "  location    - Directory that search should happen at\n")
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func displayVersion() {
	fmt.Printf("Version: %v\n", Version)
	fmt.Printf("Build:   %v\n", Build)
	fmt.Printf("Date:    %v\n", Date)
}

func setupArgs() {
	flag.Usage = Usage
	flag.BoolVar(&version, "V", false, "Version: Display version and build")
	flag.BoolVar(&commands, "C", false, "Commands: Generate archive commands")
	flag.IntVar(&days, "d", 0, "Days: Number of days to consider archiving")
	flag.StringVar(&archive, "a", "archive", "Archive: Archive path to save")
	flag.Parse()
}

func init() {
	setupArgs()

	if version {
		displayVersion()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		log.Fatal("Missing required parameters.")
	} else if len(args) > 1 {
		flag.Usage()
		log.Fatal("Too many parameters.")
	}

	location = args[0]

	info, err := os.Stat(location)
	if err != nil {
		log.Fatalf("Unable to stat location: %v", location)
	} else if info == nil {
		log.Fatalf("Location not found: %v", location)
	} else if !info.IsDir() {
		log.Fatalf("Location is not a directory: %v", location)
	}
}

func getLastModTime(location string) time.Time {
	var result time.Time
	info, err := os.Stat(location)
	if err != nil {
		log.Fatalf("Unable to stat location: %v", location)
	} else if info == nil {
		log.Fatalf("Location not found: %v", location)
	} else if info.Mode().IsRegular() {
		result = info.ModTime()
	} else if info.IsDir() {
		entries, err := os.ReadDir(location)
		if err != nil {
			log.Fatalf("Unable read dir: %v", location)
		}
		for _, entry := range entries {
			mod := getLastModTime(path.Join(location, entry.Name()))
			if mod.After(result) {
				result = mod
			}
		}
	}
	return result
}

func main() {
	var values []string
	check := time.Now().AddDate(0, 0, -days)
	entries, err := os.ReadDir(location)
	if err != nil {
		log.Fatalf("Unable to get directory of location: %v", location)
	}
	for _, entry := range entries {
		name := entry.Name()
		if name == archive {
			continue
		}
		last := getLastModTime(path.Join(location, name))
		if last.Before(check) {
			if commands {
				values = append(values, fmt.Sprintf("tar zcvf %s/%s.tar.gz %s && rm -rf %s", archive, name, name, name))
			} else {
				values = append(values, fmt.Sprintf("%s %s", last.Format("2006-01-02"), name))
			}
		}
	}
	sort.Strings(values)
	for _, value := range values {
		fmt.Println(value)
	}

}
