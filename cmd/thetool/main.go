package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

var folders = flag.String("dir", "", "The directory to load")

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		all.Usage()
		os.Exit(0)
	}

	fls := make([]*File, 0)
	for _, f := range strings.Split(*folders, ",") {
		fl, err := loadFolder(f)
		if err != nil {
			log.Fatalf("Loading folder failed: %s", err)
		}
		fls = append(fls, fl)
	}

	cmd := all.Find(args[0], args[1:])
	if err := cmd.Run(cmd, fls); err != nil {
		log.Fatalf("command failed: %s", err)
	}
}
