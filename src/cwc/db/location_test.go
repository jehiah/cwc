package db

import (
	"fmt"
	"testing"
)

func TestFindLongLatLine(t *testing.T) {
	type testCase struct {
		line      string
		lat, long float64
	}
	tests := []testCase{
		{"[ll:-78.046875,37.483577]", -78.046875, 37.483577},
		{"[ll:40.754732, -73.995206]", 40.754732, -73.995206},
		{"not ll", 0, 0},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			lat, long := findLatLongLine([]string{tc.line})
			if long != tc.long || lat != tc.lat {
				t.Errorf("got %v , %v expected %v , %v for %q", lat, long, tc.lat, tc.long, tc.line)
			}
		})
	}
}
