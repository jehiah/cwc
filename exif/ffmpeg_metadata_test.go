package exif

import (
	"fmt"
	"testing"
)

func TestParseISO6709(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		have    string
		wantLat float64
		wantLon float64
	}{
		{
			have:    "+40.7635-073.9853/",
			wantLat: 40.7635, wantLon: -73.9853,
		},
	}
	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			gotLat, gotLon, _ := parseISO6709(tc.have)
			if gotLat != tc.wantLat {
				t.Errorf("parseISO6709() gotLat = %v, want %v", gotLat, tc.wantLat)
			}
			if gotLon != tc.wantLon {
				t.Errorf("parseISO6709() gotLon = %v, want %v", gotLon, tc.wantLon)
			}
		})
	}
}
