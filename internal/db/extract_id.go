package db

import (
	"regexp"
	"strings"
)

func findServiceRequestID(lines []string) string {
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "1-1-1") && len(line) == 14:
			return "C" + line
		case strings.HasPrefix(line, "C1-1-1") && len(line) == 15:
			return line
		case strings.HasPrefix(line, "311-") && len(line) == 12:
			return line
		case strings.HasPrefix(line, "Service Request #: C1-1-1") && len(line) == 34:
			return line[19:34]
		}
	}
	return ""
}

var tlcIDPattern = regexp.MustCompile("stip #? ?(10[01][0-9]{5})s?")

func findTLCID(lines []string) string {
	for _, line := range lines {
		matches := tlcIDPattern.FindAllStringSubmatch(strings.ToLower(line), -1)
		if len(matches) >= 1 && len(matches[0]) >= 1 {
			return matches[0][1]
		}
	}
	return ""
}
