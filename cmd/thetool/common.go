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
	PartialKey  string              `json:"partial_key,omitempty" yaml:"partial_key,omitempty"`
	Key         uint32              `json:"key,omitempty" yaml:"key,omitempty"`
	Title       map[string]string   `json:"title,omitempty" yaml:"title,omitempty"`
	Description map[string]string   `json:"description,omitempty" yaml:"description,omitempty"`
	Year        int                 `json:"year,omitempty" yaml:"year,omitempty"`
	Discontinue int                 `json:"discontinue,omitempty" yaml:"discontinue,omitempty"`
	Month       int                 `json:"month" yaml:"month"`
	Day         int                 `json:"day" yaml:"day"`
	Calendar    []string            `json:"calendar,omitempty" yaml:"calendar,omitempty"`
	Holiday     map[string][]string `json:"holiday,omitempty" yaml:",omitempty"`
	Sources     []string            `json:"sources,omitempty" yaml:"sources,omitempty"`
}

func (e *Event) idx() int {
	return e.Month*100 + e.Day
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
	for i := range f.Events {
		hash := fnv.New32()
		_, _ = fmt.Fprintf(hash, "%d_%s_%d_%d_%s_%d", key0, f.Name,
			f.Events[i].Month, f.Events[i].Day, f.Events[i].PartialKey, key1)
		f.Events[i].Key = hash.Sum32()
	}
}

func openFolder(folder string) ([]string, error) {

	var ret = make([]string, 0)

	return ret, filepath.Walk(folder,
		func(path string, fi os.FileInfo, err error) error {
			if fi.IsDir() {
				return nil
			}
	
			//full := filepath.Join(folder, fi.Name())
			if ext := filepath.Ext(path); ext != ".yml" {
				return nil
			}
	
			ret = append(ret, path)
			return nil
		})
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

		for e := range f.Events {
			if f.Events[e].Key != 0 {
				return nil, fmt.Errorf("the Key should not be in the input file %q ", f.Events[e].PartialKey)
			}
		}
		res.Merge(f)
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
