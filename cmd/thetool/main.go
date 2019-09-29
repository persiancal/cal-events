package main

import (
	"flag"
	"log"
	"os"
)

var folder = flag.String("dir", "", "The directory to load")

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		all.Usage()
		os.Exit(0)
	}

	fl, err := loadFolder(*folder)
	if err != nil {
		log.Fatalf("Loading folder failed: %s", err)
	}

	cmd := all.Find(args[0], args[1:])
	if err := cmd.Run(cmd, fl); err != nil {
		log.Fatalf("command failed: %s", err)
	}
}
