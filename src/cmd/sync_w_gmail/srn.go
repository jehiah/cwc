package main

import (
	"regexp"
	"strings"
)

var sRNPattern = regexp.MustCompile("1-1-1[0-9]{9}")

func SRNFromSubject(s string) string {
	if v := sRNPattern.FindString(s); v != "" {
		return "C" + v
	}
	return ""
}

// given a message body, extracts the 311 Number from it
func SRNFromBody(lines []string) string {
	for _, line := range lines {
		if strings.HasPrefix(line, "Service Request #: C") {
			line = line[len("Service Request #: "):]
			return line[:15]
		}
	}
	return ""
}
