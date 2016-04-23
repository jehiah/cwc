package gmailutils

import (
	"encoding/base64"
	"errors"
	"log"

	"google.golang.org/api/gmail/v1"
)

// given a message ID, return it's text/body (if any)
func MessageTextBody(m *gmail.Message) ([]byte, error) {
	if m.Payload == nil {
		return nil, errors.New("no message payload")
	}
	for _, p := range m.Payload.Parts {
		if p.MimeType != "text/plain" {
			// log.Printf("skipping mime %s", p.MimeType)
			continue
		}
		decoded, err := base64.StdEncoding.DecodeString(p.Body.Data)
		return decoded, err
	}
	if len(m.Payload.Parts) == 0 && m.Payload.Body != nil {
		// m.Payload is it
		if m.Payload.MimeType != "text/plain" {
			log.Printf("skipping only mime %s", m.Payload.MimeType)
			return nil, errors.New("no text/plain found")
		}
		decoded, err := base64.StdEncoding.DecodeString(m.Payload.Body.Data)
		return decoded, err
	}
	return nil, errors.New("no text/plain found")
}

func Subject(m *gmail.Message) string {
	for _, h := range m.Payload.Headers {
		if h.Name == "Subject" {
			return h.Value
		}
	}
	return ""
}
