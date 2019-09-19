package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Event is a single event
type Event struct {
	Key         string
	Title       map[string]string
	Description map[string]string
	Year        int
	Month       int
	Day         int
	Calendar    map[string][]string
	Holiday     map[string][]string
}

// File is the single file
type File struct {
	Events []Event
}

var (
	file = flag.String("file", "", "The file to validate")
)

func isValidCountry(c string) bool {
	switch c {
	case "Iran":
		return true
	}

	return false
}

func openFile(file string) ([]byte, error) {
	fl, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fl.Close()

	return ioutil.ReadAll(fl)
}

func loadFile(data []byte) (*File, error) {
	fl := File{}
	if err := yaml.Unmarshal(data, &fl); err != nil {
		return nil, err
	}

	return &fl, nil
}

// TODO: preset support. like Persian (the current one)
func validateEventContent(ev []Event) error {
	for i := range ev {
		if ev[i].Month <= 0 || ev[i].Month > 12 {
			return fmt.Errorf("invalid month on key %d", i)
		}

		if ev[i].Month <= 6 {
			if ev[i].Day <= 0 || ev[i].Day > 31 {
				return fmt.Errorf("invalid day on key %d", i)
			}
		}

		if ev[i].Month > 6 {
			if ev[i].Day <= 0 || ev[i].Day > 30 {
				return fmt.Errorf("invalid day on key %d", i)
			}
		}

		for country := range ev[i].Holiday {
			if !isValidCountry(country) {
				return fmt.Errorf("country is invalid: %q in key %d", country, i)
			}
		}
	}

	return nil
}

func validateEventOrder(ev []Event) error {
	var (
		year, month, day int
	)

	for i := range ev {
		if ev[i].Month < month {
			return fmt.Errorf("the key %d is not in order", i)
		}

		// Reset for next month
		if ev[i].Month != month {
			day = 0
		}

		if ev[i].Day < day {
			return fmt.Errorf("the key %d is not in order", i)
		}

		if ev[i].Month == month && ev[i].Day == day {
			if ev[i].Year < year {
				return fmt.Errorf("the key %d is not in order", i)
			}
		}

		year, month, day = ev[i].Year, ev[i].Month, ev[i].Day
	}

	return nil
}

func main() {
	flag.Parse()
	data, err := openFile(*file)
	if err != nil {
		log.Fatalf("Open file failed: %s", err)
	}

	fl, err := loadFile(data)
	if err != nil {
		log.Fatalf("Loading file failed: %s", err)
	}

	if err := validateEventContent(fl.Events); err != nil {
		log.Fatalf("Validating the event data failed: %s ", err)
	}

	if err := validateEventOrder(fl.Events); err != nil {
		log.Fatalf("Validating the events order failed: %s", err)
	}
}
