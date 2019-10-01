package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEventsOrder(t *testing.T) {
	fixtures := []struct {
		events  []Event
		failKey int
	}{
		{
			events: []Event{
				{Month: 1, Day: 1, PartialKey: "PK"},
				{Month: 1, Day: 2, PartialKey: "PK"},
				{Month: 2, Day: 2, PartialKey: "PK"},
				{Month: 2, Day: 2, Year: 100, PartialKey: "PK"},
				{Month: 2, Day: 2, Year: 101, PartialKey: "PK"},
				{Month: 3, Day: 1, PartialKey: "PK"},
			},
			failKey: -1,
		},
		{
			events: []Event{
				{Month: 2, Day: 2, PartialKey: "PK"},
				{Month: 2, Day: 1, PartialKey: "PK"},
			},
			failKey: 1,
		},
	}

	for i := range fixtures {
		err := validateEventOrder(fixtures[i].events)
		if fixtures[i].failKey < 0 {
			assert.NoError(t, err)
			continue
		}

		assert.Error(t, err)
	}
}

func TestValidateEventsContent(t *testing.T) {
	fixtures := []struct {
		events  []Event
		failKey int
	}{
		{
			events: []Event{
				{Month: 0, Day: 1, PartialKey: "PK"},
			},
			failKey: 0,
		},
		{
			events: []Event{
				{Month: 2, Day: -1, PartialKey: "PK"},
			},
			failKey: 0,
		},
		{
			events: []Event{
				{Month: 7, Day: 31, PartialKey: "PK"},
			},
			failKey: 0,
		},
		{
			events: []Event{
				{Month: 7, Day: 1, Holiday: map[string][]string{"Invalid": nil}, PartialKey: "PK"},
			},
			failKey: 0,
		},
		{
			events: []Event{
				{Month: 7, Day: 1, Holiday: map[string][]string{"Iran": nil}, PartialKey: "PK"},
			},
			failKey: -1,
		},
	}

	p := &Preset{
		MonthsNormal: []int{31, 31, 31, 31, 31, 31, 30, 30, 30, 30, 30, 29},
		MonthsLeap:   []int{31, 31, 31, 31, 31, 31, 30, 30, 30, 30, 30, 30},
		MonthsName:   nil,
	}

	for i := range fixtures {
		err := validateEventContent(fixtures[i].events, p, []string{"Iran"})
		if fixtures[i].failKey < 0 {
			assert.NoError(t, err)
			continue
		}

		assert.Error(t, err)
	}

}
