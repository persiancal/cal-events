package main

import (
	"fmt"
	"os"
	"sort"

	"gopkg.in/yaml.v2"
)

var reorderOut *string

func reorder(cmd *command, fl *File) error {
	sort.Sort(fl)
	f := os.Stdout
	var err error
	if *reorderOut != "-" {
		f, err := os.Create(*reorderOut)
		if err != nil {
			return fmt.Errorf("open target file failed: %w", err)
		}
		defer func() { _ = f.Close() }()
	}

	b, err := yaml.Marshal(fl)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml: %w", err)
	}

	_, err = f.Write(b)
	return err
}

func init() {
	cmd := newCommand("reorder", "reorder the input yaml file", reorder)
	reorderOut = cmd.Flags.String("output", "-", "teh output, - for stdout")
	registerCommand(cmd)
}
