package main

import (
	"bufio"
	"bytes"
	"log"
	"strings"
	"time"
	"fmt"

	"cwc/db"
	"cwc/gmailutils"
	gmail "google.golang.org/api/gmail/v1"
)

func main() {
	var err error
	srv := gmailutils.GmailService("cwc.json")

	user := "me"
	// labels, err := gmailutils.Labels(srv, user)
	// if err != nil {
	// 	log.Fatalf("unable to fetch labels %v", err)
	// }
	// log.Printf("%#v", labels)

	// subject:"311 Service Request Closed"
	// https://godoc.org/google.golang.org/api/gmail/v1#UsersMessagesListCall
	if err != nil {
		log.Fatalf("Unable to retrieve messages. %v", err)
	}
	limit := 310
	var r *gmail.ListMessagesResponse
	for {
		// .LabelIds(labels["safe_streets/fined"])
		// .LabelIds("INBOX")
		// .Q(pattern)
		req := srv.Users.Messages.List(user).Q("subject:\"311 Service Request Closed\"").MaxResults(50)
		if r != nil {
			if r.NextPageToken != "" {
				log.Printf("getting next page %s", r.NextPageToken)
				req = req.PageToken(r.NextPageToken)
			} else {
				break
			}
		}
		r, err = req.Do()
		if err != nil {
			log.Fatalf("err getting results %s", err)
		}
		for i, m := range r.Messages {
			time.Sleep(100 * time.Millisecond)

			m, err := srv.Users.Messages.Get(user, m.Id).Do()
			if err != nil {
				log.Printf("err %s", err)
				continue
			}

			subject := gmailutils.Subject(m)
			srn := SRNFromSubject(subject)

			log.Printf("[%d] %s Subject:%s", i, srn, subject)

			if srn == "" {
				log.Printf("no 311 number found")
				continue
			}

			body, err := gmailutils.MessageTextBody(m)
			if err != nil {
				log.Printf("err %s", err)
				continue
			}
			lines := getLines(body)

			if v := SRNFromBody(lines); v != srn {
				log.Printf("missmatched SRN %q vs %q", srn, v)
				continue
			}

			complaints, err := db.Default.Find(srn[1:])
			if err != nil {
				log.Printf("%s", err)
				continue
			}
			if len(complaints) != 1 {
				log.Printf("found unexpected number of complaints %d", len(complaints))
				continue
			}
			complaint := complaints[0]

			// if we already wrote this message in... continue
			if ok, err := db.Default.ComplaintContains(complaint, m.Id); ok {
				log.Printf("already wrote message")
				continue
			} else if err != nil {
				log.Printf("%s", err)
				continue
			}

			log.Printf("** appending message to %s ***", complaint)
			ts := time.Unix(m.InternalDate / 1000, 0)
			line := interestingLine(lines)
			if strings.HasSuffix(line, "Referred to S") {
				// fix a 311 bug where text emails are truncated at '&'
				line += "&E"
			}
			
			logMsg := fmt.Sprintf("\n[email:%s %s] %s\n", m.Id, ts.Format("2006/01/02 15:04"), line)
			err = db.Default.Append(complaint, logMsg)
			if err != nil {
				log.Printf("%s", err)
			}
			
			limit--
			if limit <= 0 {
				log.Printf("at limit. ending")
				return
			}
			// if i == 0 {
			// 	log.Printf("body %s", body)
			// }
		}

	}
}

func getLines(b []byte) []string {
	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("%s", err)
	}
	return lines
}

// the last "useful" line is the one before 'Get Service Request Details'
func interestingLine(lines []string) string {
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "The TLC is investigating") {
			return lines[i]
		}
		if lines[i] == "Get Service Request Details" && i > 0 {
			return lines[i-1]
		}
	}
	return ""
}
