package main

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jehiah/cwc/internal/db"
	"github.com/jehiah/cwc/internal/gmailutils"
	gmail "google.golang.org/api/gmail/v1"
)

type NoticeOfDecision struct {
	DB             db.ReadWrite
	ArchiveMessage MessageArchiver
	*gmail.UsersMessagesAttachmentsService
}

func (s *NoticeOfDecision) BuildQuery(u *gmail.UsersMessagesListCall) *gmail.UsersMessagesListCall {
	return u.LabelIds("INBOX").Q("subject:\"notice of decision\" from:@tlc.nyc.gov")
}

func (s *NoticeOfDecision) Handle(m *gmail.Message) error {
	prettyID := prettyMessageID(m)

	subject := gmailutils.Subject(m)
	fields := strings.Fields(subject)
	var TLCComplaintNumber string

	if strings.HasSuffix(fields[len(fields)-1], "c") {
		TLCComplaintNumber = fields[len(fields)-1]
		TLCComplaintNumber = TLCComplaintNumber[:len(TLCComplaintNumber)-1]
	}

	log.Printf("%s TLC Complaint:%s Subject:%s", prettyID, TLCComplaintNumber, subject)

	if TLCComplaintNumber == "" {
		log.Printf("no TLC request number found")
		// fmt.Printf("%s", body)
		return nil
	}

	complaints, err := s.DB.Find(TLCComplaintNumber)
	if err != nil {
		log.Printf("%s", err)
		return nil
	}

	if len(complaints) != 1 {
		log.Printf("found unexpected number of complaints %d", len(complaints))
		return nil
	}
	complaint := complaints[0]
	log.Printf("\tfound related complaint %s", complaint)

	// find PDF
	attachmentID, err := gmailutils.MessagePDF(m)
	if err != nil {
		log.Printf("no PDF found %s", err)
		return err
	}
	log.Printf("\tdownloading attachment %s", attachmentID)

	msgPart, err := s.UsersMessagesAttachmentsService.Get("me", m.Id, attachmentID).Do()
	if err != nil {
		log.Printf("error getting attachment %s", err)
		return err
	}
	pdfBody, err := base64.StdEncoding.DecodeString(msgPart.Data)
	if err != nil {
		pdfBody, err = base64.URLEncoding.DecodeString(msgPart.Data)
	}
	if err != nil {
		log.Printf("error base64 decoding attachment %s", err)
		return err
	}

	log.Printf("\tPDF size %d", len(pdfBody))

	fullPath := filepath.Join(s.DB.FullPath(complaint), TLCComplaintNumber+"c.pdf")
	log.Printf("\tcreating %s", fullPath)

	f, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil && os.IsExist(err) {
		log.Printf("\tfile already exists")
		return nil
	}
	if err != nil {
		log.Printf("%s", err)
		return err
	}
	f.Write(pdfBody)
	f.Close()

	log.Printf("\tarchiving email")
	return s.ArchiveMessage(m.Id)
}
