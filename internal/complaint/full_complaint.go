package complaint

import (
	"bytes"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jehiah/cwc/internal/reg"
)

type FullComplaint struct {
	Complaint   Complaint `json:"complaint"`
	Timestamp   int64     `json:"timestamp"`
	Time        time.Time `json:"-"`
	License     string    `json:"license_plate"`
	VehicleType string    `json:"vehicle_type"`
	Location    string    `json:"location"`
	Address     string    `json:"address"`
	Description string    `json:"description"`

	Status           State     `json:"status"`
	ServiceRequestID string    `json:"311_service_request_id"`
	TLCID            string    `json:"tlc_id,omitempty"`
	Hearing          time.Time `json:"hearing"`
	Violations       []reg.Reg `json:"violations"`
	Tweets           []string  `json:"tweets,omitempty"`

	Body   string   `json:"body"`
	Lines  []string `json:"lines"` // the non-empty lines of text
	Photos []string `json:"photos"`
	Videos []string `json:"videos"`
	PDFs   []string `json:"pdfs"`
	Files  []string `json:"files"`
	Lat    float64  `json:"lat,omitempty"`
	Long   float64  `json:"long,omitempty"`

	PhotoDetails []Photo `json:"photo_details"`
}

func (fc FullComplaint) IsNewSRNumberFormat() bool {
	return strings.HasPrefix(fc.ServiceRequestID, "311-")
}

var tweetPattern = regexp.MustCompile(`https?://[^\s]+`)

func ParseComplaint(c RawComplaint, files []string) (*FullComplaint, error) {
	b := string(c.Body)
	contains := func(pattern string) bool {
		return bytes.Contains(c.Body, []byte(pattern))
	}
	lines := splitLines(b)
	f := &FullComplaint{
		Complaint:   c.Complaint,
		Timestamp:   c.Time().Unix(),
		Time:        c.Time(),
		License:     c.License(),
		VehicleType: reg.Taxi.String(),

		Body:             b,
		Lines:            lines,
		Hearing:          extractHearingDate(lines),
		Status:           DetectState(b),
		TLCID:            findTLCID(lines),
		ServiceRequestID: findServiceRequestID(lines),
		Tweets:           tweetPattern.FindAllString(b, -1),
	}
	f.Lat, f.Long = findLatLongLine(lines)

	if contains("FHV") {
		f.VehicleType = reg.FHV.String()
	}
	for _, r := range reg.All {
		if contains(r.Code) {
			f.Violations = append(f.Violations, r)
		}
	}
	if len(f.Lines) >= 1 {
		location := strings.SplitN(f.Lines[0], " ", 5)
		if len(location) == 5 {
			f.Location = location[4]
		}
	}
	if len(f.Lines) >= 2 {
		f.Description = f.Lines[1]
		for _, line := range lines[1:] {
			if strings.HasPrefix(line, "Address: ") {
				f.Address = strings.TrimSpace(line[len("Address: "):])
			}
		}
	}

	for _, filename := range files {
		ext := strings.ToLower(filepath.Ext(filename))
		switch ext {
		case ".mov", ".m4v", ".mp4":
			f.Videos = append(f.Videos, filename)
		case ".bmp", ".png", ".jpg", ".jpeg", ".tif", ".gif", ".heic":
			f.Photos = append(f.Photos, filename)
		case ".pdf":
			f.PDFs = append(f.PDFs, filename)
			if filename == f.TLCID+"c.pdf" {
				f.Status = NoticeOfDecision
			}
		default:
			f.Files = append(f.Files, filename)
		}
	}

	// if Unknown and >6month, set to expired
	if f.Status == Unknown && time.Since(f.Time) > (time.Hour*4320) {
		f.Status = Expired
	}

	return f, nil
}

// splitLines returns the non-empty trimmed lines
func splitLines(s string) []string {
	var o []string
	for _, l := range strings.Split(s, "\n") {
		line := strings.TrimSpace(l)
		if line != "" {
			o = append(o, line)
		}
	}
	return o
}

// // fullComplaintsByAge implements sort.Interface for []Complaint based on
// // the complaint Time.
// type FullComplaintsByAge []*FullComplaint
// func (a FullComplaintsByAge) Len() int           { return len(a) }
// func (a FullComplaintsByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// func (a FullComplaintsByAge) Less(i, j int) bool { return a[i].Time.Before(a[j].Time) }

// FullComplaintsByHearing implements sort.Interface for []Complaint based on
// the complaint Hearing if exists otherwise Time. Hearing always sorts before those without a hearing
type FullComplaintsByHearing []*FullComplaint

func (a FullComplaintsByHearing) Len() int      { return len(a) }
func (a FullComplaintsByHearing) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a FullComplaintsByHearing) Less(i, j int) bool {
	switch {
	case !a[i].Hearing.IsZero() && !a[j].Hearing.IsZero():
		return a[i].Hearing.Before(a[j].Hearing)
	case a[i].Hearing.IsZero() && a[j].Hearing.IsZero():
		return a[i].Time.Before(a[j].Time)
	case a[i].Hearing.IsZero():
		return true
	case a[j].Hearing.IsZero():
		return false
	}
	panic("unreachable")
}

func (f FullComplaint) Contains(pattern string) bool {
	return strings.Contains(f.Body, pattern)
}
