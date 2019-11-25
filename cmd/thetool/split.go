package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	target *string
	base   *string
)

func monthBase(fls []*File) error {
	for _, fl := range fls {
		ev := make([][]Event, len(fl.Months.Normal))
		for i := range fl.Events {
			fl.Events[i].Key = 0
			ev[fl.Events[i].Month-1] = append(ev[fl.Events[i].Month-1], fl.Events[i])
		}

		for i := range ev {
			name := strings.ToLower(filepath.Join(*target, fmt.Sprintf("%02d-%s.yml", i+1, fl.Months.Name[i]["en_US"])))
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
	}

	return nil
}

func eventBase(fls []*File) error {
	for _, fl := range fls {
		sort.Sort(fl)
		evs := make(map[string][]Event)
		for i := range fl.Events {
			fl.Events[i].Key = 0
			k := fmt.Sprintf("%02d%02d", fl.Events[i].Month, fl.Events[i].Day)
			if _, ok := evs[k]; ok {
				evs[k] = append(evs[k], fl.Events[i])
			} else {
				evs[k] = []Event{fl.Events[i]}
			}
		}

		

		for k, v := range evs {
			name := strings.ToLower(filepath.Join(*target,
				fl.Name,
				fmt.Sprintf("%02d-%s", v[0].Month, fl.Months.Name[v[0].Month-1]["en_US"]),
				fmt.Sprintf("%s.yml", k[2:])))
			b, err := yaml.Marshal(File{
				Events: v,
			})
			if err != nil {
				return err
			}
			os.MkdirAll(strings.ToLower(filepath.Join(*target,
				fl.Name,
				fmt.Sprintf("%02d-%s", v[0].Month, fl.Months.Name[v[0].Month-1]["en_US"]))), os.ModePerm)

			if err := ioutil.WriteFile(name, b, 0600); err != nil {
				return err
			}
		}

		fl.Events = nil
		name := strings.ToLower(filepath.Join(*target,fl.Name, "preset.yml"))
		b, err := yaml.Marshal(fl)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(name, b, 0600); err != nil {
			return err
		}
	}

	return nil
}

func split(cmd *command, fls []*File) error {
	switch *base {
	case "month":
		return monthBase(fls)
	case "event":
		return eventBase(fls)
	default:
		return fmt.Errorf("the %q is not valid value for -base", *base)
	}
}

func init() {
	cur, _ := os.Getwd()
	cmd := newCommand("split", "split one year based on the months", split)
	target = cmd.Flags.String("dist", cur, "The dist folder")
	base = cmd.Flags.String("base", cur, `The base of split. valid values are "month" and "event" `)
	registerCommand(cmd)
}
