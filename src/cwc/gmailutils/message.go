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
	var gotError error
	for _, p := range part.Parts {
		b, err := recurse(p)
		if b != nil {
			return b, err
		}
		if err != nil {
			gotError = err
		}
	}
	log.Printf("%s", part.MimeType)
	switch {
	case strings.HasPrefix(part.MimeType, "text/plain"):
		b, err := base64.StdEncoding.DecodeString(part.Body.Data)
		if err != nil {
			return base64.URLEncoding.DecodeString(part.Body.Data)
		}
		return b, err
	}
	return nil, gotError
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
