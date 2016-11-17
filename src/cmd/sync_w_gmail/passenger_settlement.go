package main

import (
	"fmt"
	"log"

	"cwc/db"
	"cwc/gmailutils"

	gmail "google.golang.org/api/gmail/v1"
)

type PassengerSettlementNotification struct {
	DB db.DB
}

func (s *PassengerSettlementNotification) findSRN(lines []string) string {
	l := FirstLineWithPrefix("Subject: TLC Complaint # 1-1", lines, false)
	if l != "" {
		return l[len("Subject: TLC Complaint # "):][:14]
	}
	return l
}

func (s *PassengerSettlementNotification) BuildQuery(u *gmail.UsersMessagesListCall) *gmail.UsersMessagesListCall {
	return u.LabelIds("INBOX").Q("subject:\"Passenger Settlement Notification\"")
}

func (s *PassengerSettlementNotification) Handle(m *gmail.Message) error {
	prettyID := prettyMessageID(m)

	body, err := gmailutils.MessageTextBody(m)
	if err != nil {
		log.Printf("err %s", err)
		return nil
	}
	lines := getLines(body)
	srn := s.findSRN(lines)
	log.Printf("%s %s Subject:%s", prettyID, srn, gmailutils.Subject(m))

	if srn == "" {
		log.Printf("no 311 service request number found")
		return nil
	}

	complaints, err := db.Default.Find(srn)
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

	line := "The driver has pleaded guilty to a rule violation and has paid a penalty."
	log.Print(line)
	err = s.DB.Append(complaint, fmt.Sprintf("\n%s %s\n", prettyID, line))
	return err
}
