package main

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	partialValidator = regexp.MustCompile("^[a-z0-9_]+$")
	multipleSpace    = regexp.MustCompile(`\s{2,}`)
)

func reference(a ...interface{}) string {
	r := ""
	if len(a) == 0 {
		return r
	}
	for i := range a {
		r += fmt.Sprintf("%v => ", a[i])
	}
	return r[:len(r)-3]
}

func textValidator(s string, rs ...interface{}) error {
	r := reference(rs...)
	if strings.TrimSpace(s) != s {
		return fmt.Errorf("validate text failed: start and/or end with space in %q", r)
	}
	if multipleSpace.MatchString(s) {
		return fmt.Errorf("validate text failed: double or more white spaces in %q", r)
	}
	return nil
}

func validate(cmd *command, fl *File) error {
	if err := validateEventCalendar(fl); err != nil {
		return fmt.Errorf("validate calendars failed: %w", err)
	}

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

func validateEventContent(ev []Event, p *Months, countries []string) error {
	for i := range ev {
		if ev[i].Key != 0 {
			return fmt.Errorf("the Key should not be in the input file")
		}

		for l, t := range ev[i].Title {
			if err := textValidator(t, ev[i].PartialKey, "title", l); err != nil {
				return err
			}
		}

		for l, t := range ev[i].Description {
			if err := textValidator(t, ev[i].PartialKey, "description", l); err != nil {
				return err
			}
		}

		for l, t := range ev[i].Calendar {
			if err := textValidator(t, ev[i].PartialKey, "calendar", l, t); err != nil {
				return err
			}
		}

		for l, t := range ev[i].Holiday {
			for p, r := range t {
				if err := textValidator(r, ev[i].PartialKey, "holiday", l, p); err != nil {
					return err
				}
			}
		}

		if !partialValidator.MatchString(ev[i].PartialKey) {
			return fmt.Errorf("the partial key %q is invalid, only lower english chars, _ and numbers are allowed ([a-z0-9_])", ev[i].PartialKey)
		}

		if ev[i].Month <= 0 || ev[i].Month > len(p.Normal) {
			return fmt.Errorf("invalid month on key %d", i)
		}

		max := p.Normal[ev[i].Month-1]
		if leap := p.Leap[ev[i].Month-1]; leap > max {
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
				return fmt.Errorf("the key %q is not in order", ev[i].PartialKey)
			}
		}

		year, lastIdx = ev[i].Year, ev[i].idx()
	}

	return nil
}

func validateEventCalendar(fl *File) error {
	for _, ev := range fl.Events {
		if len(ev.Calendar) == 0 {
			return fmt.Errorf("event %q has no calendar", ev.PartialKey)
		}
	middleLoop:
		for _, c := range ev.Calendar {
			for _, cv := range fl.Calendars {
				if cv["en_US"] == c {
					continue middleLoop
				}
			}
			return fmt.Errorf("calendar %q for event %q is invalid", c, ev.PartialKey)
		}
	}

	return nil
}

func init() {
	cmd := newCommand("validate", "validate the input yaml file", validate)
	registerCommand(cmd)
}
