package shared

import (
	"testing"
	"time"

	"github.com/timenglesf/personal-site/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 8, 17, 0, 0, 0, 0, time.UTC),
			want: "August 17, 2024",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := HumanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}
}
