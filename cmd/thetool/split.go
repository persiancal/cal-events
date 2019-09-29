package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var target *string

func split(cmd *command, fl *File) error {
	ev := make([][]Event, len(fl.Months.MonthsNormal))
	for i := range fl.Events {
		ev[fl.Events[i].Month-1] = append(ev[fl.Events[i].Month-1], fl.Events[i])
	}

	for i := range ev {
		name := strings.ToLower(filepath.Join(*target, fmt.Sprintf("%02d-%s.yml", i+1, fl.Months.MonthsName[i]["en_US"])))
		b, err := yaml.Marshal(File{
			Events: ev[i],
		})
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(name, b, 0600); err != nil {
			return err
		}
	}

	fl.Events = nil
	name := filepath.Join(*target, "preset.yml")
	b, err := yaml.Marshal(fl)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(name, b, 0600); err != nil {
		return err
	}

	return nil
}

func init() {
	cur, _ := os.Getwd()
	cmd := newCommand("split", "split one year based on the months", split)
	target = cmd.Flags.String("dist", cur, "The dist folder")
	registerCommand(cmd)
}
