package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"cwc/db"
	"cwc/gmailutils"

	"golang.org/x/net/context"
	gmail "google.golang.org/api/gmail/v1"
)

func main() {
	limit := flag.Int("limit", 500, "max number to process")
	flag.Parse()

	srv := gmailutils.GmailService("cwc.json")
	bg := context.TODO()

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
		&SettlementNotification{DB: db.Default, alternate: true},
		&SettlementNotification{DB: db.Default},
		// &ServiceReqeustUpdate{db.Default},
	}
	for _, h := range handlers {
		q := h.BuildQuery(srv.Users.Messages.List(user)).MaxResults(50)
		var c int
		err := q.Pages(bg, func(r *gmail.ListMessagesResponse) error {
			for _, m := range r.Messages {
				c++
				if c > *limit {
					log.Printf("over max %d", *limit)
					return fmt.Errorf("handled %d (over max %d)", c, *limit)
				}
				time.Sleep(100 * time.Millisecond)
				var err error
				m, err = srv.Users.Messages.Get(user, m.Id).Do()
				if err != nil {
					log.Printf("%s", err)
					return err
				}
				err = h.Handle(m)
				if err != nil {
					log.Printf("%s", err)
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Fatalf("%s", err)
		}
	}
}