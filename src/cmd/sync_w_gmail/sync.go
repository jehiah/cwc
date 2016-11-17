package main

import (
	"fmt"
	"log"
	"time"

	"cwc/db"
	"cwc/gmailutils"
	"golang.org/x/net/context"
	gmail "google.golang.org/api/gmail/v1"
)

func main() {
	srv := gmailutils.GmailService("cwc.json")
	bg := context.Background()

	user := "me"
	// labels, err := gmailutils.Labels(srv, user)
	// if err != nil {
	// 	log.Fatalf("unable to fetch labels %v", err)
	// }
	// log.Printf("%#v", labels)

	// subject:"311 Service Request Closed"
	// https://godoc.org/google.golang.org/api/gmail/v1#UsersMessagesListCall
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve messages. %v", err)
	// }
	handlers := []EmailHandler{
		&ServiceReqeustUpdate{db.Default},
	}
	for _, h := range handlers {
		max := 500
		q := h.BuildQuery(srv.Users.Messages.List(user)).MaxResults(50)
		var c int
		err := q.Pages(bg, func(r *gmail.ListMessagesResponse) error {
			for _, m := range r.Messages {
				c++
				if c > max {
					return fmt.Errorf("handled %d (over max %d)", c, max)
				}
				time.Sleep(100 * time.Millisecond)
				var err error
				m, err = srv.Users.Messages.Get(user, m.Id).Do()
				if err != nil {
					return err
				}
				err = h.Handle(m)
				if err != nil {
					log.Printf("%s", err)
					return err
				}
				return nil
			}
			return nil
		})
		if err != nil {
			log.Fatalf("%s", err)
		}
	}
}
