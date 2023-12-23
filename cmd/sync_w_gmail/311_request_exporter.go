package main

import (
	"fmt"
	"log"

	"github.com/jehiah/cwc/internal/gmailutils"
	gmail "google.golang.org/api/gmail/v1"
)

type NYC311RequestExporter struct {
	version int
}

func (s *NYC311RequestExporter) BuildQuery(u *gmail.UsersMessagesListCall) *gmail.UsersMessagesListCall {
	switch s.version {
	case 1:
		// SRNotification@customerservice.nyc.gov
		// Confirmation | Updated | Closed
		// 311 Service Request Closed #: C1-1-1626059981 , Street Sign - Missing
	default:
		// label == nyc/311
		return u.LabelIds("Label_7662191922466997049").Q("from:SRNotice@customercare.nyc.gov subject:\"SR Submitted #\" after:2023-10-07")
	}
	return nil
}

func (s *NYC311RequestExporter) Handle(m *gmail.Message) error {
	prettyID := prettyMessageID(m)

	subject := gmailutils.Subject(m)
	srn := SRNFromSubject(subject)

	if srn == "" {
		log.Printf("warning: no 311 service request number found")
		return nil
	}

	fmt.Printf("%s %s\n", srn, prettyID)
	return nil
}
