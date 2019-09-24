package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Event is a single event
type Event struct {
	Key         string              `json:"key" yaml:"key"`
	Title       map[string]string   `json:"title,omitempty" yaml:"title,omitempty"`
	Description map[string]string   `json:"description,omitempty" yaml:"description,omitempty"`
	Year        int                 `json:"year,omitempty" yaml:"year,omitempty"`
	Month       int                 `json:"month" yaml:"month"`
	Day         int                 `json:"day" yaml:"day"`
	Calendar    map[string][]string `json:"calendar" yaml:"calendar"`
	Holiday     map[string][]string `json:"holiday,omitempty" yaml:",omitempty"`
}

func (e *Event) idx() int {
	return e.Month*12 + e.Day
}

// File is the single file
type File struct {
	Name   string  `json:"-" yaml:"-"`
	Events []Event `json:"events" yaml:"events"`
}

func openFile(file string) ([]byte, error) {
	fl, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fl.Close()

	return ioutil.ReadAll(fl)
}

func loadFile(name string, data []byte) (*File, error) {
	fl := File{
		Name: name,
	}
	if err := yaml.Unmarshal(data, &fl); err != nil {
		return nil, err
	}

	return &fl, nil
}
