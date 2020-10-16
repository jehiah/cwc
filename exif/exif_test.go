package exif

import (
	"reflect"
	"testing"
	"time"
)

func mustParse(t time.Time, err error) time.Time {
	if err != nil {
		panic(err.Error())
	}
	return t
}

func TestParse(t *testing.T) {
	type testCase struct {
		filename string
		expect   *Exif
	}
	tests := []testCase{
		{"testdata/IMG_4056.JPG", &Exif{
			Lat:          40.75985277777778,
			Long:         -73.99134722222222,
			Created:      mustParse(time.Parse("2006-01-02 15:04:05 -0700 MST", "2016-11-04 17:28:19 -0400 EDT")),
			ExifRotation: 90,
		}},
	}
	for _, tc := range tests {
		got, err := Parse(tc.filename)
		if err != nil {
			t.Fatalf("got err %s", err)
		}
		if !reflect.DeepEqual(got, tc.expect) {
			t.Logf("ts %v", got.Created)
			t.Errorf("got %#v expected %#v", got, tc.expect)
		}
	}
}
