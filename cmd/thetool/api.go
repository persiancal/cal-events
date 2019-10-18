package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type categoryHandler interface {
	Group(fl *File, ret map[string][]*Event) map[string][]*Event
}

type dailyCat struct {
	root string
}

func (dc dailyCat) path(name string, month, day int) string {
	return filepath.Join(dc.root, strings.ToLower(name), fmt.Sprint(month), fmt.Sprint(day))
}

func (dc dailyCat) Group(fl *File, ret map[string][]*Event) map[string][]*Event {
	for m := range fl.Months.Name {
		for d := 0; d < fl.Months.Leap[m] || d < fl.Months.Normal[m]; d++ {
			path := dc.path(fl.Name, m+1, d+1)
			if _, ok := ret[path]; !ok {
				ret[path] = []*Event{}
			}
		}
	}

	for i := range fl.Events {
		path := dc.path(fl.Name, fl.Events[i].Month, fl.Events[i].Day)
		ret[path] = append(ret[path], &fl.Events[i])
	}

	return ret
}

func allHandlers(root string) []categoryHandler {
	return []categoryHandler{
		dailyCat{root: root},
	}
}

func categorize(fl *File, handlers ...categoryHandler) map[string][]*Event {
	ret := make(map[string][]*Event)
	for i := range handlers {
		ret = handlers[i].Group(fl, ret)
	}
	return ret
}

func writeFiles(m map[string][]*Event) error {
	for i := range m {
		b, err := json.MarshalIndent(m[i], "", "  ")
		if err != nil {
			return fmt.Errorf("convert to json failed: %w", err)
		}

		path := i + ".json"
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0740); err != nil {
			return fmt.Errorf("creating directory %q failed: %w", dir, err)
		}

		if err := ioutil.WriteFile(path, b, 0600); err != nil {
			return fmt.Errorf("write file %q failed: %w", path, err)
		}
	}

	return nil
}

func writeStaticApi(fls []*File, root string) error {
	handlers := allHandlers(root)
	for i := range fls {
		m := categorize(fls[i], handlers...)
		if err := writeFiles(m); err != nil {
			return err
		}
	}

	return nil
}
