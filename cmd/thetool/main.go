package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var file = flag.String("file", "", "base yaml file")

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		all.Usage()
		os.Exit(0)
	}
	data, err := openFile(*file)
	if err != nil {
		log.Fatalf("Open file failed: %s", err)
	}

	name := strings.TrimSuffix(filepath.Base(*file), filepath.Ext(*file))
	fl, err := loadFile(name, data)
	if err != nil {
		log.Fatalf("Loading file failed: %s", err)
	}

	cmd := all.Find(args[0], args[1:])
	if err := cmd.Run(cmd, fl); err != nil {
		log.Fatalf("Command failed: %s", err)
	}
}
