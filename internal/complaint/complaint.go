package complaint

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func init() {
	geoclientCache = make(map[string]LL)

}

func New(dt time.Time, license string) Complaint {
	c := fmt.Sprintf("%s_%s", dt.Format("20060102_1504"), license)
	return Complaint(c)
}

type Complaint string

func (c Complaint) String() string {
	return fmt.Sprintf("%s - %s", c.License(), c.Time().Format("Mon Jan 2 2006 3:04pm"))
}
func (c Complaint) ID() string {
	return string(c)
}

type RawComplaint struct {
	Complaint
	Body []byte
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
type ComplaintsByAge []Complaint

func (a ComplaintsByAge) Len() int           { return len(a) }
func (a ComplaintsByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ComplaintsByAge) Less(i, j int) bool { return a[i].Time().Before(a[j].Time()) }
