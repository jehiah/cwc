package main

import (
	"fmt"
	"log"
	"strings"

	"cwc/db"
	"cwc/gmailutils"

	gmail "google.golang.org/api/gmail/v1"
)

type ServiceReqeustUpdate struct {
	DB db.DB
	force bool
	dryrun bool
}

func (s *ServiceReqeustUpdate) BuildQuery(u *gmail.UsersMessagesListCall) *gmail.UsersMessagesListCall {
	return u.LabelIds("INBOX").Q("subject:\"311 Service Request Update\" OR subject:\"311 Service Request Closed\"")
}

func (s *ServiceReqeustUpdate) Handle(m *gmail.Message) error {
	prettyID := prettyMessageID(m)

	subject := gmailutils.Subject(m)
	srn := SRNFromSubject(subject)

	log.Printf("%s %s Subject:%s", prettyID, srn, subject)

	if srn == "" {
		log.Printf("warning: no 311 service request number found")
		return nil
	}

	body, err := gmailutils.MessageTextBody(m)
	if err != nil {
		log.Printf("error %s", err)
		return nil
	}
	lines := getLines(body)

	if v := SRNFromBody(lines); v != srn {
		log.Printf("error: missmatched SRN %q vs %q", srn, v)
		return nil
	}

	complaints, err := db.Default.Find(srn[1:])
	if err != nil {
		log.Printf("%s", err)
		return nil
	}
	if len(complaints) != 1 {
		log.Printf("error: found unexpected number of complaints %d %v", len(complaints), complaints)
		return nil
	}
	complaint := complaints[0]

	log.Printf("\tfound related complaint %s", complaint)
	// if we already wrote this message in... continue
	if ok, err := db.Default.ComplaintContains(complaint, m.Id); ok {
		log.Printf("\talready recorded")
		if s.force == false {
			return nil
		}
	} else if err != nil {
		log.Printf("%s", err)
		return nil
	}

	line := interestingLine(lines)
	if strings.HasSuffix(line, "Referred to S") {
		// fix a 311 bug where text emails are truncated at '&'
		line += "&E"
	}

	// strip these prefixes from the line
	// clearly this should be a regex of sorts
	for _, s := range []string{
		"The TLC is investigating your complaint and will contact youif further information is needed.",
		"The TLC is investigating your complaint and will contact you if further information is needed.",
		"The TLC is investigating your complaint and will contact youif further information is needed",
		"The TLC is investigating your complaint and will contact you if further information is needed",
		"The TLC is investigating your complaint and will contact ou if further information is needed.",
		"The TLC is investigating youur complaint and will contact you if further information is needed.",
		"The TLC is investigating your complaint and will contact if further information is needed.",
		"TheTLC is investigating your complaint and will contact you if further information is needed.",
		"The TLC is Investigating your complaint and will contact you if further information is needed.",
		"The TLC is investigating your complint and will contact you if further information is needed.",
		"The TLC is investigating your complint and will contact you if further information is needed.",
	} {
		if strings.HasSuffix(line, s) {
			log.Printf("\tskipping line %q", line)
			return nil
		}
		if strings.HasPrefix(line, s) {
			line = strings.TrimSpace(line[len(s):])
		}
	}
	if line == "" {
		return nil
	}

	log.Printf("\t> %s", line)
	if !s.dryrun {
		err = s.DB.Append(complaint, fmt.Sprintf("\n%s %s\n", prettyID, line))
	}
	return err
}

// the last "useful" line is the one before 'Get Service Request Details'
func interestingLine(lines []string) string {
	// line before 'Get Service Request Details' is most likely interesting
	for i, line := range lines {
		if line != "Get Service Request Details" {
			continue
		}
		if i > 1 && lines[i-2] == "The Taxi and Limousine Commission was unable to identify the driver or car service company named in your complaint." {
			return lines[i-1] + " " + lines[i-2]
		}
		if i > 0 {
			return lines[i-1]
		}
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "The TLC is investigating") {
			return line
		}
	}
	return strings.Join(lines, " ")
}
