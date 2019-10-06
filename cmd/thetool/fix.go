package main

import (
	"fmt"
)

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func fixCalendar(fl *File, p []map[string]string) error {
	for i := range fl.Events {
		ev := &fl.Events[i]
		ca, ok := ev.Calendar["en_US"]
		if !ok {
			return fmt.Errorf("event calendar is empty: %q", ev.PartialKey)
		}

	bigLoop:
		for _, c := range ca {
			if c == "Iran" {
				ev.NewCalendar = append(ev.NewCalendar, "Iran Official")
				continue
			}
			for _, key := range p {
				cc := key["en_US"]
				if c == cc {
					ev.NewCalendar = append(ev.NewCalendar, cc)
					continue bigLoop
				}
			}
			return fmt.Errorf("invalid calendar %q for event %q", c, ev.PartialKey)
		}
		ev.NewCalendar = unique(ev.NewCalendar)
	}

	return nil
}
