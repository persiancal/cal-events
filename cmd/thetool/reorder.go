package main

import (
	"fmt"
	"os"
	"sort"

	"gopkg.in/yaml.v2"
)

// TODO : this is a utility, and works with one file only as target, fix it or remove it

var reorderOut *string

func reorder(cmd *command, fls []*File) error {
	for _, fl := range fls {

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
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	cmd := newCommand("reorder", "reorder the input yaml file", reorder)
	reorderOut = cmd.Flags.String("output", "-", "teh output, - for stdout")
	registerCommand(cmd)
}
