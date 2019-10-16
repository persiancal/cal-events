package main

import "fmt"

func unique(_ *command, fls []*File) error {
	keys := make(map[uint32]Event)
	for _, fl := range fls {
		for i := range fl.Events {
			if d, ok := keys[fl.Events[i].Key]; ok {
				return fmt.Errorf("duplicate key, the key %d is same for both partial keys %q and %q", fl.Events[i].Key, d.PartialKey, fl.Events[i].PartialKey)
			}
			keys[fl.Events[i].Key] = fl.Events[i]
		}
	}
	return nil
}

func init() {
	cmd := newCommand("unique", "make sure the event key is unique across all events", unique)
	registerCommand(cmd)
}
