package main

import (
	"fmt"
	"strings"
)

func validate(cmd *command, fl *File) error {
	if err := validateEventContent(fl.Events, fl.Months, fl.Countries); err != nil {
		return fmt.Errorf("validate events failed: %w", err)
	}

	if err := validateEventOrder(fl.Events); err != nil {
		return fmt.Errorf("validating the events order failed: %w", err)
	}

	return nil
}

func isValidCountry(c string, list []string) bool {
	for i := range list {
		if list[i] == c {
			return true
		}
	}

	return false
}

func validateEventContent(ev []Event, p *Preset, countries []string) error {
	for i := range ev {
		if ev[i].Key != 0 {
			return fmt.Errorf("the Key should not be in the input file")
		}

		if strings.Trim(ev[i].PartialKey, "\n\t ") == "" {
			return fmt.Errorf("the partial key is empty")
		}

		if ev[i].Month <= 0 || ev[i].Month > len(p.MonthsNormal) {
			return fmt.Errorf("invalid month on key %d", i)
		}

		max := p.MonthsNormal[ev[i].Month-1]
		if leap := p.MonthsLeap[ev[i].Month-1]; leap > max {
			max = leap
		}

		if ev[i].Day <= 0 || ev[i].Day > max {
			return fmt.Errorf("invalid day on key %d", i)
		}

		for country := range ev[i].Holiday {
			if !isValidCountry(country, countries) {
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
			return fmt.Errorf("the key %d is not in order %+v => %+v", i, ev[i], ev[i-1])
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
