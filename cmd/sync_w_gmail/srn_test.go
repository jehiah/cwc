package main

import (
	"fmt"
	"testing"
	"time"
)

func TestSRNFromSubject(t *testing.T) {
	type testCase struct {
		subject, expect string
	}
	tests := []testCase{
		{"SR Updated # 311-09751295", "311-09751295"},
		{"SR Submitted # 311-09751295", "311-09751295"},
	}
	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			if got := SRNFromSubject(tc.subject); got != tc.expect {
				t.Errorf("got %q expected %q for %q", got, tc.expect, tc.subject)
			}
		})
	}
}

func TestTLCIDFromSubject(t *testing.T) {
	type testCase struct {
		subject, expect string
	}

	tests := []testCase{
		{"notice  of adjournment 10091665c", "10091665"},
		{"tlc notice of adjournment 10092127c", "10092127"},
		{"REVISED ******* notice of adjournment 10091414c", "10091414"},
		{"motion to vacate 10090018c", "10090018"},
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
		{
			"   New Hearing Date:                       December 21, 2017",
			"Time:              9:30 AM",
			true,
			time.Date(2017, 12, 21, 9, 30, 0, 0, time.UTC),
		},
		{
			"A hearing on this summons will take place at 31-00 47th Ave, 3rd Floor, Long Island City, NY 11101, on 4/24/2023 at 10:30 AM.", "",
			true,
			time.Date(2023,4,24,10,30,0,0, time.UTC),
		},

	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			lines := []string{tc.line1, tc.line2}
			got, ok := HearingDateFromBody(lines)
			if ok != tc.expect {
				t.Fatalf("got %v - %v expected %v for %q %q", got, ok, tc.expect, tc.line1, tc.line2)
			}
			if !got.Equal(tc.time) {
				t.Errorf("got %s expected %s for %q %q", got, tc.time, tc.line1, tc.line2)
			}
		})
	}

}
