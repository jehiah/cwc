package main

import (
	gmail "google.golang.org/api/gmail/v1"
)

type EmailHandler interface {
	BuildQuery(*gmail.UsersMessagesListCall) *gmail.UsersMessagesListCall
	Handle(*gmail.Message) error
}
