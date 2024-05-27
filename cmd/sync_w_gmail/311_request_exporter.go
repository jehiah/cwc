package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jehiah/cwc/internal/gmailutils"
	"github.com/jehiah/cwc/internal/nycapi"
	gmail "google.golang.org/api/gmail/v1"
)

type NYC311RequestExporter struct {
	version int
	nycAPI  nycapi.Client
}

func (s *NYC311RequestExporter) BuildQuery(u *gmail.UsersMessagesListCall) *gmail.UsersMessagesListCall {
	switch s.version {
	case 1:
		// SRNotification@customerservice.nyc.gov
		// Confirmation | Updated | Closed
		// 311 Service Request Closed #: C1-1-1626059981 , Street Sign - Missing
	default:
		// label == nyc/311
		return u.LabelIds("Label_7662191922466997049").Q("from:SRNotice@customercare.nyc.gov subject:\"SR Submitted #\" after:2023-12-31")
	}
	return nil
}

func (s *NYC311RequestExporter) Handle(m *gmail.Message) error {
	prettyID := prettyMessageID(m)

	subject := gmailutils.Subject(m)
	srn := SRNFromSubject(subject)

	if srn == "" {
		log.Printf("warning: no 311 service request number found %s", prettyID)
		return nil
	}

	srAPI, err := s.nycAPI.GetServiceRequest(context.Background(), srn)
	if err != nil {
		log.Printf("error fetching 311 service request %s: %v", srn, err)
		time.Sleep(time.Second * 2)
		return nil
	}
	body, _ := json.Marshal(srAPI)
	fmt.Printf("%s\n", body)
	time.Sleep(500 * time.Millisecond)
	return nil
}
