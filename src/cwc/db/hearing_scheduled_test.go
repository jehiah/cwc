package db

import (
	"testing"
	"time"
	"fmt"
)

func TestExtractHearing(t *testing.T) {
	type testCase struct {
		line     string
		expected time.Time
	}
	ny, _ := time.LoadLocation("America/New_York")
	tests := []testCase{
		{"NT. Hearing scheduled 10/28/16 at 10:00 AM. ar", time.Date(2016, 10, 28, 10, 0, 0, 0, ny)},
		{"hm Hearing scheduled 5/6/2016 at 11:00 AM", time.Date(2016, 5, 6, 11, 0, 0, 0, ny)},
		{"hm  HEaring scheduled 4/5/16 10:00am JW2", time.Date(2016, 4, 5, 10, 0, 0, 0, ny)},
		{"mailed to driver on 5/09/16 5/09/16 Hearing Scheduled 7/5/16 at 2:30 PM", time.Date(2016, 7, 5, 14, 30, 0, 0, ny)},
		{"Hearing scheduled on 8/21/2017 at 9:00 am.", time.Date(2017,8,21,9,0,0,0, ny)},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t * testing.T){
			got := extractHearingDate([]string{tc.line})
			if !got.Equal(tc.expected) {
				t.Errorf("got %q expected %q for %q", got, tc.expected, tc.line)
			}
		})
	}
}
