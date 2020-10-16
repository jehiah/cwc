package main

import (
	"bufio"
	"bytes"
	"log"
	"regexp"
	"strings"
	"time"
)

var sRNPattern = regexp.MustCompile("1-1-1[0-9]{9}")
var adjournmentPattern = regexp.MustCompile(`NOTICE\s+OF\s+(?:HEARING|ADJOURNMENT)\s+(?:AND HEARING LOCATION\s+)?(10[0-9]{6})C`)
var vacatePattern = regexp.MustCompile(`MOTION TO VACATE\s+(10[0-9]{6})C`)

func SRNFromSubject(s string) string {
	if v := sRNPattern.FindString(s); v != "" {
		return "C" + v
	}
	return ""
}

// TLCIDFromSubject - parses "notice of XYZ 1234c"
func TLCIDFromSubject(s string) string {
	s = strings.ToUpper(s)
	if v := adjournmentPattern.FindStringSubmatch(s); len(v) >= 2 {
		return v[1]
	}
	if v := vacatePattern.FindStringSubmatch(s); len(v) >= 2 {
		return v[1]
	}
	return ""
}

func TLCIDFromBody(lines []string) string {
	for _, line := range lines {
		if strings.HasPrefix(line, "Re: Complaint No.: ") && strings.HasSuffix(line, "C") {
			f := strings.Fields(line)
			return f[len(f)-1]
		}
	}
	return ""
}

func HearingDateFromBody(lines []string) (time.Time, bool) {
	var dateStr, timeStr string
	for _, line := range lines {
		if idx := strings.Index(line, "The new date is "); idx != -1 {
			dateStr = strings.Trim(line[idx+16:], ".")
		}
		if idx := strings.Index(line, "The new time is "); idx != -1 {
			timeStr = strings.Trim(line[idx+16:], ".")
		}
		if idx := strings.Index(line, "New Hearing Date:"); idx != -1 {
			dateStr = strings.Trim(line[idx+17:], ".")
			dateStr = strings.TrimSpace(dateStr)
		}
		if idx := strings.Index(line, "Time: "); idx != -1 && (strings.HasSuffix(line, " AM") || strings.HasSuffix(line, " PM")) {
			timeStr = strings.TrimSpace(strings.Trim(line[idx+5:], "."))
		}
	}
	if dateStr == "" || timeStr == "" {
		return time.Time{}, false
	}
	//  December 1, 2017 2:30 PM
	t, err := time.Parse("January 2, 2006 3:04 PM", dateStr+" "+timeStr)
	if err != nil {
		log.Printf("\t%s %s", dateStr, timeStr)
		return time.Time{}, false
	}
	return t, true
}

// given a message body, extracts the 311 Number from it
func SRNFromBody(lines []string) string {
	s := FirstLineWithPrefix("Service Request #: C", lines, true)
	if s != "" {
		return ("C" + s)[:15]
	}
	return ""
}

func SRNFromTLCComplaintBody(lines []string) string {
	type matcher struct {
		pattern string
		prefix  string
	}
	for _, m := range []matcher{
		{"Subject: TLC Complaint # 1-1", "Subject: TLC Complaint # "},
		{"Subject: TLC Complaint #1-1", "Subject: TLC Complaint #"},
		{"Subject: TLC Complaint 1-1", "Subject: TLC Complaint "},
	} {
		if line := FirstLineWithPrefix(m.pattern, lines, false); line != "" {
			return line[len(m.prefix):]
		}
	}
	return ""
}

func FirstLineWithPrefix(needle string, lines []string, strip bool) string {
	for _, line := range lines {
		if strings.HasPrefix(line, needle) {
			if strip {
				return line[len(needle):]
			}
			return line
		}
	}
	return ""
}

func getLines(b []byte) []string {
	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("%s", err)
	}
	return lines
}
