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
				{Month: 1, Day: 1},
				{Month: 1, Day: 2},
				{Month: 2, Day: 2},
				{Month: 2, Day: 2, Year: 100},
				{Month: 2, Day: 2, Year: 101},
				{Month: 3, Day: 1},
			},
			failKey: -1,
		},
		{
			events: []Event{
				{Month: 2, Day: 2},
				{Month: 2, Day: 1},
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
				{Month: 0, Day: 1},
			},
			failKey: 0,
		},
		{
			events: []Event{
				{Month: 2, Day: -1},
			},
			failKey: 0,
		},
		{
			events: []Event{
				{Month: 7, Day: 31},
			},
			failKey: 0,
		},
	}

	for i := range fixtures {
		err := validateEventContent(fixtures[i].events)
		if fixtures[i].failKey < 0 {
			assert.NoError(t, err)
			continue
		}

		assert.Error(t, err)
	}

}
