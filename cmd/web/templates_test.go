package main

import (
	"testing"
	"time"

	"github.com/thedevscott/trug/internal/assert"
)

func TestBasicHumanDate(t *testing.T) {
	tm := time.Date(2026, 3, 17, 10, 15, 0, 0, time.UTC)
	hd := humanDate(tm)

	if hd != "17 Mar 2026" {
		t.Errorf("got %q; want %q", hd, "17 Mar 2026")
	}
}

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2026, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2026",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)
		})
	}
}
