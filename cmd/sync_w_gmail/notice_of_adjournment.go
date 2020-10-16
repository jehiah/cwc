package main

import (
	"fmt"
	"log"

	"github.com/jehiah/cwc/internal/db"
	"github.com/jehiah/cwc/internal/gmailutils"
	gmail "google.golang.org/api/gmail/v1"
)

type NoticeOfAdjournment struct {
	DB             db.DB
	ArchiveMessage MessageArchiver
	Alternate      bool
}

func (s *NoticeOfAdjournment) BuildQuery(u *gmail.UsersMessagesListCall) *gmail.UsersMessagesListCall {
	if s.Alternate {
		return u.LabelIds("INBOX").Q("from:@tlc.nyc.gov to:tlccomplaintunit@tlc.nyc.gov subject:\"motion to vacate\"")
	}
	return u.LabelIds("INBOX").Q("from:@tlc.nyc.gov to:tlccomplaintunit@tlc.nyc.gov subject:\"notice of adjournment\"")
}

func (s *NoticeOfAdjournment) Handle(m *gmail.Message) error {
	prettyID := prettyMessageID(m)

	body, err := gmailutils.MessageTextBody(m)
	if err != nil {
		log.Printf("err %s", err)
		return nil
	}

	subject := gmailutils.Subject(m)
	tlcid := TLCIDFromSubject(subject)
	log.Printf("%s %s Subject:%s", prettyID, tlcid, subject)

	if tlcid == "" {
		log.Printf("no complaint number found in subject")
	}

	lines := getLines(body)

	if tlcid == "" {
		tlcid = TLCIDFromBody(lines)
	}

	if tlcid == "" {
		log.Printf("no complaint number found in body")
		return nil
	}

	hearing, ok := HearingDateFromBody(lines)
	if !ok {
		log.Printf("no hearing found")
		return nil
	}
	log.Printf("\tHearing scheduled for %s", hearing)

	complaints, err := db.Default.Find(tlcid)
	if err != nil {
		log.Printf("%s", err)
		return nil
	}
	if len(complaints) != 1 {
		log.Printf("found unexpected number of complaints %d for %s", len(complaints), tlcid)
		return nil
	}
	complaint := complaints[0]
	log.Printf("\tfound related complaint %s", complaint)

	// if we already wrote this message in... continue
	if ok, err := db.Default.ComplaintContains(complaint, m.Id); ok {
		log.Printf("\talready recorded")
		log.Printf("\tarchiving email")
		return s.ArchiveMessage(m.Id)
	} else if err != nil {
		log.Printf("%s", err)
		return nil
	}
	line := "NOTICE OF ADJOURNMENT"
	if s.Alternate {
		line = "MOTION TO VACATE"
	}
	line += " hearing scheduled for " + hearing.Format("01/02/06 at 3:04 PM")
	log.Printf("\t> %s", line)
	err = s.DB.Append(complaint, fmt.Sprintf("\n%s %s\n", prettyID, line))
	if err != nil {
		return err
	}
	log.Printf("\tarchiving email")
	return s.ArchiveMessage(m.Id)
	return nil
}
