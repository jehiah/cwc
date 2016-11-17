package gmailutils

import (
	"encoding/base64"
	"errors"
	"log"
	"strings"

	"google.golang.org/api/gmail/v1"
)

func recurse(part *gmail.MessagePart) ([]byte, error) {
	if part == nil || part.Body == nil {
		return nil, nil
	}
	for _, p := range part.Parts {
		b, err := recurse(p)
		if err != nil || b != nil {
			return b, err
		}
	}
	log.Printf("%s", part.MimeType)
	switch {
	case strings.HasPrefix(part.MimeType, "text/plain"):
		 return base64.StdEncoding.DecodeString(part.Body.Data)
	}
	return nil, nil
}

// given a message ID, return it's text/body (if any)
func MessageTextBody(m *gmail.Message) ([]byte, error) {
	body, err := recurse(m.Payload)
	if body == nil {
		return nil, errors.New("no message payload")
	}
	return body, err
}

func Subject(m *gmail.Message) string {
	for _, h := range m.Payload.Headers {
		if h.Name == "Subject" {
			return h.Value
		}
	}
	return ""
}
