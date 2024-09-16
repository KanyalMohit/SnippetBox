package main

import (
	"testing"
	"time"

	"snippetbox.mohit.net/internal/assert"
)

func TestHumanData(t *testing.T) {

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 8, 22, 22, 5, 0, 0, time.UTC),
			want: "22 Aug 2024 at 22:05",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 8, 22, 22, 5, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "22 Aug 2024 at 22:05",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)
			
		})
	}

}
