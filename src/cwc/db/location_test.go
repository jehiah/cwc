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

func TestRandomBetween(t *testing.T) {
	type testCase struct {
		A, B LL
	}
	tests := []testCase{
		{LL{40.76311, -73.999}, LL{40.762155, -73.997407}},
		{LL{40.760790, -73.998308}, LL{40.762757, -73.996956}},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			extent := func(a, b float64) (float64, float64) {
				if a > b {
					return b, a
				}
				return a, b
			}
			got := randomBetween(tc.A, tc.B)
			if min, max := extent(tc.A.Lat, tc.B.Lat); min > got.Lat || max < got.Lat {
				t.Errorf("Lat %#v not between %#v and %#v", got, tc.A, tc.B)
			}
			if min, max := extent(tc.A.Long, tc.B.Long); min > got.Long || max < got.Long {
				t.Errorf("Long %#v not between %#v and %#v", got, tc.A, tc.B)
			}
		})
	}

}
