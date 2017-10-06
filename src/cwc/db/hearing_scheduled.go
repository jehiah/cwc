package db

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var nyc *time.Location

func init() {
	nyc, _ = time.LoadLocation("America/New_York")
}

var hearingPattern = regexp.MustCompile("(Scheduled|scheduled|hearing sch) (?:on |for )?([0-9]{1,2}/[0-9]{1,2}/2?0?1[0-9]) (at |- )?([0-9]{1,2}:?[0-9]{0,2}) ?(am|AM|pm|PM)")
var hearingLayouts = []string{
	"1/2/06 3:04pm",
	"1/2/2006 3:04pm",
}

func extractHearingDate(lines []string) time.Time {
	// iterate backwards start with last log line
	for i := len(lines); i > 0; i-- {
		line := lines[i-1]
		matches := hearingPattern.FindAllStringSubmatch(line, -1)
		if len(matches) >= 1 && len(matches[0]) >= 6 {
			s := fmt.Sprintf("%s %s%s", matches[0][2], matches[0][4], strings.ToLower(matches[0][5]))
			for _, layout := range hearingLayouts {
				t, err := time.ParseInLocation(layout, s, nyc)
				if err == nil {
					return t
				}
			}
		}
	}

	return time.Time{}
}
