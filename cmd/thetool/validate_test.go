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
				{Month: 1, Day: 1, PartialKey: "partial_key"},
				{Month: 1, Day: 2, PartialKey: "partial_key"},
				{Month: 2, Day: 2, PartialKey: "partial_key"},
				{Month: 2, Day: 2, Year: 100, PartialKey: "partial_key"},
				{Month: 2, Day: 2, Year: 101, PartialKey: "partial_key"},
				{Month: 3, Day: 1, PartialKey: "partial_key"},
			},
			failKey: -1,
		},
		{
			events: []Event{
				{Month: 2, Day: 2, PartialKey: "partial_key"},
				{Month: 2, Day: 1, PartialKey: "partial_key"},
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
				{Month: 0, Day: 1, PartialKey: "partial_key"},
			},
			failKey: 0,
		},
		{
			events: []Event{
				{Month: 2, Day: -1, PartialKey: "partial_key"},
			},
			failKey: 0,
		},
		{
			events: []Event{
				{Month: 7, Day: 31, PartialKey: "partial_key"},
			},
			failKey: 0,
		},
		{
			events: []Event{
				{Month: 7, Day: 1, Holiday: map[string][]string{"Invalid": nil}, PartialKey: "partial_key"},
			},
			failKey: 0,
		},
		{
			events: []Event{
				{Month: 7, Day: 1, Holiday: map[string][]string{"Iran": nil}, PartialKey: "partial_key"},
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

func TestTextValidator(t *testing.T) {
	type args struct {
		s string
		r string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "should return multiple space error",
			args:    args{s: "space  space", r: "reference"},
			wantErr: true,
		},
		{
			name:    "should return trim error",
			args:    args{s: "trim it ", r: "reference"},
			wantErr: true,
		},
		{
			name:    "should return trim error",
			args:    args{s: " trim it", r: "reference"},
			wantErr: true,
		},
		{
			name:    "should not return error",
			args:    args{s: "Correct format", r: "reference"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := textValidator(tt.args.s,tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("textValidator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
