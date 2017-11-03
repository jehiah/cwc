package main

import (
	"fmt"
	"testing"
	"time"
)

func TestTLCIDFromSubject(t *testing.T) {
	type testCase struct {
		subject, expect string
	}

	tests := []testCase{
		{"notice  of adjournment 10091665c", "10091665"},
		{"tlc notice of adjournment 10092127c", "10092127"},
		{"REVISED ******* notice of adjournment 10091414c", "10091414"},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			if got := TLCIDFromSubject(tc.subject); got != tc.expect {
				t.Errorf("got %q expected %q for %q", got, tc.expect, tc.subject)
			}
		})
	}
}

func TestHearingDateFromBody(t *testing.T) {
	type testCase struct {
		line1, line2 string
		expect       bool
		time         time.Time
	}

	tests := []testCase{
		{
			"X                         The new date is December 1, 2017.",
			"______                   Time Changed. The new time is 2:30 PM.",
			true,
			time.Date(2017, 12, 1, 14, 30, 0, 0, time.UTC),
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			lines := []string{tc.line1, tc.line2}
			got, ok := HearingDateFromBody(lines)
			if ok != tc.expect {
				t.Fatalf("got %v expected %v for %q %q", ok, tc.expect, tc.line1, tc.line2)
			}
			if !got.Equal(tc.time) {
				t.Errorf("got %s expected %s for %q %q", got, tc.time, tc.line1, tc.line2)
			}
		})
	}

}
