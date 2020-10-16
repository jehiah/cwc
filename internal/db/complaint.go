package db

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Complaint string

func (c Complaint) String() string {
	return fmt.Sprintf("%s - %s", c.License(), c.Time().Format("Mon Jan 2 2006 3:04pm"))
}
func (c Complaint) ID() string {
	return string(c)
}

func (c Complaint) Time() time.Time {
	if len(c) < 13 {
		log.Panicf("invalid format %s", string(c))
	}
	t, err := time.Parse("20060102_1504", string(c)[:13])
	if err != nil {
		log.Printf("err parsing time %s", err)
		return time.Time{}
	}
	return t
}

func (c Complaint) License() string {
	chunks := strings.SplitN(string(c), "_", 3)
	if len(chunks) == 3 {
		return chunks[2]
	}
	return ""
}

// complaintsByAge implements sort.Interface for []Complaint based on
// the Time function.
type complaintsByAge []Complaint

func (a complaintsByAge) Len() int           { return len(a) }
func (a complaintsByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a complaintsByAge) Less(i, j int) bool { return a[i].Time().Before(a[j].Time()) }
