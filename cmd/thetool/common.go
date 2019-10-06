package main

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Rest in peace...
const (
	key0 = 0x61_73_68_6b_61_6e
	key1 = 0x76_61_6e_65_68
)

// Event is a single event
type Event struct {
	PartialKey  string            `json:"partial_key,omitempty" yaml:"partial_key,omitempty"`
	Key         uint32            `json:"key,omitempty" yaml:"key,omitempty"`
	Title       map[string]string `json:"title,omitempty" yaml:"title,omitempty"`
	Description map[string]string `json:"description,omitempty" yaml:"description,omitempty"`
	Year        int               `json:"year,omitempty" yaml:"year,omitempty"`
	Month       int               `json:"month" yaml:"month"`
	Day         int               `json:"day" yaml:"day"`
	// deprecated, transitional phase
	Calendar    map[string][]string `json:"calendar" yaml:"calendar"`
	NewCalendar []string            `json:"new_calendar,omitempty" yaml:"new_calendar,omitempty"`
	Holiday     map[string][]string `json:"holiday,omitempty" yaml:",omitempty"`
	Sources     []string            `json:"sources,omitempty" yaml:"sources,omitempty"`
}

func (e *Event) idx() int {
	return e.Month*100 + e.Day
}

func (e *Event) CalculateKey(collection string) {
	hash := fnv.New32()
	_, _ = fmt.Fprintf(hash, "%d_%s_%d_%d_%s_%d", key0, collection, e.Month, e.Day, e.PartialKey, key1)
	e.Key = hash.Sum32()
}

// Months is the month structure validator
type Months struct {
	Normal []int               `json:"normal,omitempty" yaml:"normal,omitempty"`
	Leap   []int               `json:"leap,omitempty" yaml:"leap,omitempty"`
	Name   []map[string]string `json:"name,omitempty" yaml:"name,omitempty"`
}

// File is the single file
type File struct {
	Name      string              `json:"name,omitempty" yaml:"name,omitempty"`
	Countries []string            `json:"countries,omitempty" yaml:"countries,omitempty"`
	Calendars []map[string]string `json:"calendars,omitempty" yaml:"calendars,omitempty"`
	Months    *Months             `json:"months,omitempty" yaml:"months,omitempty"`
	Events    []Event             `json:"events,omitempty" yaml:"events,omitempty"`
}

func (f *File) Len() int {
	return len(f.Events)
}

func (f *File) Less(i, j int) bool {
	if f.Events[i].idx() == f.Events[j].idx() {
		return f.Events[i].Year < f.Events[j].Year
	}

	return f.Events[i].idx() < f.Events[j].idx()
}

func (f *File) Swap(i, j int) {
	f.Events[i], f.Events[j] = f.Events[j], f.Events[i]
}

func (f *File) Merge(new *File) {
	// TODO: Some file should only have events, and one only have preset, validate them
	if f.Name == "" {
		f.Name = new.Name
	}
	if f.Months == nil {
		f.Months = new.Months
	}

	if len(f.Calendars) == 0 {
		f.Calendars = new.Calendars
	}

	f.Countries = append(f.Countries, new.Countries...)
	f.Events = append(f.Events, new.Events...)
}

func openFolder(folder string) ([]string, error) {
	fl, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	var ret = make([]string, 0, len(fl))
	for i := range fl {
		if fl[i].IsDir() {
			continue
		}

		full := filepath.Join(folder, fl[i].Name())
		if ext := filepath.Ext(full); ext != ".yml" {
			continue
		}

		ret = append(ret, full)
	}

	return ret, nil
}

func loadFolder(folder string) (*File, error) {
	fl, err := openFolder(folder)
	if err != nil {
		return nil, err
	}

	res := &File{}

	for i := range fl {
		data, err := openFile(fl[i])
		if err != nil {
			return nil, err
		}

		f, err := loadFile(data)
		if err != nil {
			return nil, err
		}
		res.Merge(f)
	}

	for i := range fl {
		data, err := openFile(fl[i])
		if err != nil {
			return nil, err
		}

		f, err := loadFile(data)
		if err != nil {
			return nil, err
		}

		if err := fixCalendar(f, res.Calendars); err != nil {
			return nil, err
		}
		d, err := yaml.Marshal(f)
		if err != nil {
			panic(err)
		}

		if err := ioutil.WriteFile(fl[i], d, 0600); err != nil {
			panic(err)
		}
	}

	return res, nil
}

func openFile(file string) ([]byte, error) {
	fl, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() { _ = fl.Close() }()

	return ioutil.ReadAll(fl)
}

func loadFile(data []byte) (*File, error) {
	fl := File{}
	if err := yaml.Unmarshal(data, &fl); err != nil {
		return nil, err
	}

	return &fl, nil
}
