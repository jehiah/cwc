package db

import (
	"regexp"
	"time"
	"log"
	"strings"
	"fmt"
)

var ny *time.Location
func init() {
	ny, _ = time.LoadLocation("America/New_York")
}
var hearingPattern = regexp.MustCompile("(s|S)cheduled ([0-9]{1,2}/[0-9]{1,2}/2?0?1[0-9]) (at )?([0-9]{1,2}:?[0-9]{0,2}) ?(am|AM|pm|PM)")
var hearingLayouts = []string{
	"1/2/06 3:04pm",
	"1/2/2006 3:04pm",
}
func extractHearingDate(lines []string) time.Time {
	for _, line := range lines {
		matches := hearingPattern.FindAllStringSubmatch(line, -1)
		if len(matches) >= 1 && len(matches[0]) >= 6{
			s := fmt.Sprintf("%s %s%s", matches[0][2], matches[0][4], strings.ToLower(matches[0][5]))
			log.Printf("match %#v = %q", matches[0], s)
			for _, layout := range hearingLayouts {
				t, err := time.ParseInLocation(layout, s, ny)
				if err == nil {
					return t
				} else {
					log.Printf("%q error %s", s, err)
				}
			}
		}
	}
	
	return time.Time{}
}
