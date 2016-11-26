package main

import (
	"bufio"
	"bytes"
	"log"
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
