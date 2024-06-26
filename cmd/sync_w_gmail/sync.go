package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jehiah/cwc/db"
	"github.com/jehiah/cwc/internal/gmailutils"
	"github.com/jehiah/cwc/internal/nycapi"
	"golang.org/x/net/context"
	gmail "google.golang.org/api/gmail/v1"
)

func main() {
	limit := flag.Int("limit", 100, "max number to process")
	dryRun := flag.Bool("dry-run", false, "dry run")
	export := flag.Bool("export", false, "export 311 numbers")
	flag.Parse()

	srv := gmailutils.GmailService("cwc.json")
	bg := context.TODO()

	user := "me"

	var archiveMessage = func(id string) error {
		if *dryRun {
			return nil
		}
		call := srv.Users.Messages.Modify(user, id, &gmail.ModifyMessageRequest{RemoveLabelIds: []string{"INBOX"}})
		_, err := call.Do()
		return err
	}
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
	attachmentSvc := gmail.NewUsersMessagesAttachmentsService(srv)
	handlers := []EmailHandler{
		&ServiceReqeustUpdate{DB: db.Default, ArchiveMessage: archiveMessage},
		&ServiceReqeustUpdate{DB: db.Default, ArchiveMessage: archiveMessage, version: 1},
		&NoticeOfHearing{DB: db.Default, ArchiveMessage: archiveMessage},
		&NoticeOfAdjournment{DB: db.Default, ArchiveMessage: archiveMessage},
		&NoticeOfAdjournment{DB: db.Default, ArchiveMessage: archiveMessage, Alternate: true},
		&NoticeOfDecision{DB: db.Default, ArchiveMessage: archiveMessage, UsersMessagesAttachmentsService: attachmentSvc},
		&SettlementNotification{DB: db.Default, version: 1, ArchiveMessage: archiveMessage},
		&SettlementNotification{DB: db.Default, version: 2, ArchiveMessage: archiveMessage},
		&SettlementNotification{DB: db.Default, ArchiveMessage: archiveMessage},
	}
	if *export {
		handlers = []EmailHandler{
			&NYC311RequestExporter{
				nycAPI: nycapi.Client{
					SubscriptionKey: os.Getenv("OCP_APM_KEY"),
				},
			},
		}
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
