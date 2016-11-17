package main

import (
	"fmt"
	"time"

	gmail "google.golang.org/api/gmail/v1"
)

type EmailHandler interface {
	BuildQuery(*gmail.UsersMessagesListCall) *gmail.UsersMessagesListCall
	Handle(*gmail.Message) error
}

func prettyMessageID(m *gmail.Message) string {
	ts := time.Unix(m.InternalDate/1000, 0)
	return fmt.Sprintf("[email:%s %s]", m.Id, ts.Format("2006/01/02 15:04"))
}
