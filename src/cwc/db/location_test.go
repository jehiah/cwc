package db

import (
	"fmt"
	"testing"
)

func TestFindLongLatLine(t *testing.T) {
	type testCase struct {
		line      string
		long, lat float64
	}
	tests := []testCase{
		{"[ll:-78.046875,37.483577]", -78.046875, 37.483577},
		{"not ll", 0, 0},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			long, lat := findLongLatLine([]string{tc.line})
			if long != tc.long || lat != tc.lat {
				t.Errorf("got %v , %v expected %v , %v for %q", long, lat, tc.long, tc.lat, tc.line)
			}
		})
	}
}
