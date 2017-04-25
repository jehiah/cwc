package gmailutils

import (
	"encoding/base64"
	"errors"
	"strings"

	"google.golang.org/api/gmail/v1"
)

func recurse(part *gmail.MessagePart, mimeType string) ([]byte, string, error) {
	if part == nil || part.Body == nil {
		return nil, "", nil
	}
	var gotError error
	for _, p := range part.Parts {
		b, m, err := recurse(p, mimeType)
		if b != nil || m != "" {
			return b, m, err
		}
		if err != nil {
			gotError = err
		}
	}
	switch {
	case strings.HasPrefix(part.MimeType, mimeType):
		if part.Body.AttachmentId != "" {
			return nil, part.Body.AttachmentId, nil
		}
		b, err := base64.StdEncoding.DecodeString(part.Body.Data)
		if err != nil {
			b, err = base64.URLEncoding.DecodeString(part.Body.Data)
		}
		return b, "", err
	}
	return nil, "", gotError
}

// given a message ID, return it's text/body (if any)
func MessageTextBody(m *gmail.Message) ([]byte, error) {
	body, _, err := recurse(m.Payload, "text/plain")
	if body == nil {
		return nil, errors.New("no message payload")
	}
	return body, err
}

func MessagePDF(m *gmail.Message) (attachmentID string, err error) {
	_, attachmentID, err = recurse(m.Payload, "application/pdf")
	if attachmentID == "" {
		return "", errors.New("no application/pdf payload")
	}
	return
}

func Subject(m *gmail.Message) string {
	for _, h := range m.Payload.Headers {
		if h.Name == "Subject" {
			return h.Value
		}
	}
	return ""
}
