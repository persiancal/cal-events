package main

import (
	"fmt"
)

func validate(cmd *command, fl *File) error {
	if err := validateEventContent(fl.Events); err != nil {
		return fmt.Errorf("validate events failed: %w", err)
	}

	if err := validateEventOrder(fl.Events); err != nil {
		return fmt.Errorf("validating the events order failed: %w", err)
	}

	return nil
}

func isValidCountry(c string) bool {
	switch c {
	case "Iran":
		return true
	}

	return false
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
		lastIdx, year int
	)

	for i := range ev {
		if lastIdx > ev[i].idx() {
			return fmt.Errorf("the key %d is not in order", i)
		}

		if lastIdx == ev[i].idx() {
			if ev[i].Year < year {
				return fmt.Errorf("the key %d is not in order", i)
			}
		}

		year, lastIdx = ev[i].Year, ev[i].idx()
	}

	return nil
}

func init() {
	cmd := newCommand("validate", "validate the input yaml file", validate)
	registerCommand(cmd)
}
