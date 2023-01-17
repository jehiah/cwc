package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jehiah/cwc/internal/db"
	"github.com/jehiah/cwc/internal/gmailutils"
	gmail "google.golang.org/api/gmail/v1"
)

type SettlementNotification struct {
	DB             db.ReadWrite
	version        int
	ArchiveMessage MessageArchiver
}

func (s *SettlementNotification) BuildQuery(u *gmail.UsersMessagesListCall) *gmail.UsersMessagesListCall {
	if s.version == 1 {
		return u.LabelIds("INBOX").Q("from:@tlc.nyc.gov to:tlccomplaintunit@tlc.nyc.gov \"No hearing will be necessary as the driver has plead guilty to an appropriate charge and paid a penalty.\"")
	}
	if s.version == 2 {
		return u.LabelIds("INBOX").Q("from:@tlc.nyc.gov to:tlccomplaintunit@tlc.nyc.gov \"The driver has plead guilty to a rule violation and has paid the penalty. The case has now closed and there will be no hearing.\"")
	}

	return u.LabelIds("INBOX").Q("subject:\"Passenger Settlement Notification\"")
}

func (s *SettlementNotification) Handle(m *gmail.Message) error {
	prettyID := prettyMessageID(m)

	body, err := gmailutils.MessageTextBody(m)
	if err != nil {
		log.Printf("err %s", err)
		return nil
	}
	lines := getLines(body)
	srn := SRNFromTLCComplaintBody(lines)
	var TLCComplaintNumber string
	if srn != "" {
		if srn != "" && strings.Contains(srn, "/") {
			v := strings.SplitN(srn, "/", 2)
			srn, TLCComplaintNumber = strings.TrimSpace(v[0]), strings.TrimSpace(v[1])
		}
	}

	log.Printf("%s %s Subject:%s", prettyID, srn, gmailutils.Subject(m))

	if srn == "" {
		log.Printf("no 311 service request number found")
		// fmt.Printf("%s", body)
		return nil
	}

	complaints, err := db.Default.Find(srn)
	if err != nil {
		log.Printf("%s", err)
		return nil
	}
	if len(complaints) != 1 {
		log.Printf("found unexpected number of complaints %d for %q", len(complaints), srn)
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

	line := "The driver has pleaded guilty to a rule violation and has paid a penalty."
	if TLCComplaintNumber != "" {
		line = fmt.Sprintf("complaint %s. %s", TLCComplaintNumber, line)
	}
	log.Printf("\t> %s", line)
	err = s.DB.Append(complaint, fmt.Sprintf("\n%s %s\n", prettyID, line))
	if err != nil {
		return err
	}
	log.Printf("\tarchiving email")
	return s.ArchiveMessage(m.Id)
}
