package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"cwc/db"
	"cwc/gmailutils"
	gmail "google.golang.org/api/gmail/v1"
)

type ServiceReqeustUpdate struct {
	DB db.DB
}

func (s *ServiceReqeustUpdate) BuildQuery(u *gmail.UsersMessagesListCall) *gmail.UsersMessagesListCall {
	return u.LabelIds("INBOX").Q("subject:\"311 Service Request Update\" OR subject:\"311 Service Request Closed\"")
}

func (s *ServiceReqeustUpdate) Handle(m *gmail.Message) error {
	ts := time.Unix(m.InternalDate/1000, 0)
	prettyID := fmt.Sprintf("[email:%s %s]", m.Id, ts.Format("2006/01/02 15:04"))

	subject := gmailutils.Subject(m)
	srn := SRNFromSubject(subject)

	log.Printf("%s %s Subject:%s", prettyID, srn, subject)

	if srn == "" {
		log.Printf("no 311 number found")
		return nil
	}

	body, err := gmailutils.MessageTextBody(m)
	if err != nil {
		log.Printf("err %s", err)
		return nil
	}
	lines := getLines(body)

	if v := SRNFromBody(lines); v != srn {
		log.Printf("missmatched SRN %q vs %q", srn, v)
		return nil
	}

	complaints, err := db.Default.Find(srn[1:])
	if err != nil {
		log.Printf("%s", err)
		return nil
	}
	if len(complaints) != 1 {
		log.Printf("found unexpected number of complaints %d", len(complaints))
		return nil
	}
	complaint := complaints[0]

	// if we already wrote this message in... continue
	if ok, err := db.Default.ComplaintContains(complaint, m.Id); ok {
		log.Printf("already recorded - %s", complaint)
		return nil
	} else if err != nil {
		log.Printf("%s", err)
		return nil
	}

	log.Printf("** message related to %s ***", complaint)
	line := interestingLine(lines)
	if strings.HasSuffix(line, "Referred to S") {
		// fix a 311 bug where text emails are truncated at '&'
		line += "&E"
	}

	if strings.HasSuffix(line, "will contact you if further information is needed.") {
		log.Printf("skipping %q", line)
		return nil
	}

	log.Print(line)
	err = s.DB.Append(complaint, fmt.Sprintf("\n%s %s\n", prettyID, line))
	return err
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

// the last "useful" line is the one before 'Get Service Request Details'
func interestingLine(lines []string) string {
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "The TLC is investigating") {
			return lines[i]
		}
		if lines[i] == "Get Service Request Details" && i > 0 {
			return lines[i-1]
		}
	}
	return ""
}
