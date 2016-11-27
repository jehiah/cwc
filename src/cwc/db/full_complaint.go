package db

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"cwc/reg"
)

type FullComplaint struct {
	Complaint   Complaint `json:"complaint"`
	Timestamp   int64     `json:"timestamp"`
	Time        time.Time `json:"-"`
	License     string    `json:"license_plate"`
	VehicleType string    `json:"vehicle_type"`
	Location    string    `json:"location"`
	Description string    `json:"description"`

	Status           State     `json:"status"`
	ServiceRequestID string    `json:"311_service_request_id"`
	TLCID            string    `json:"tlc_id,omitempty"`
	Hearing          bool      `json:"hearing"`
	Violations       []reg.Reg `json:"violations"`
	Tweets           []string  `json:"tweets,omitempty"`

	Body     string   `json:"body"`
	Lines    []string `json:"lines"` // the non-empty lines of text
	BasePath string   `json:"-"`
	Photos   []string `json:"photos"`
	Videos   []string `json:"videos"`
	Files    []string `json:"files"`

	PhotoDetails []*Photo `json:"photo_details"`
}

func (d DB) FullComplaint(c Complaint) (*FullComplaint, error) {
	f, err := d.Open(c)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	body, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	path := d.FullPath(c)
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	files, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	return ParseComplaint(c, body, path, files)
}

var tweetPattern = regexp.MustCompile(`https?://[^\s]+`)

func ParseComplaint(c Complaint, body []byte, path string, files []string) (*FullComplaint, error) {
	b := string(body)
	contains := func(pattern string) bool {
		return bytes.Contains(body, []byte(pattern))
	}
	f := &FullComplaint{
		Complaint:   c,
		Timestamp:   c.Time().Unix(),
		Time:        c.Time(),
		License:     c.License(),
		VehicleType: reg.Taxi.String(),

		Body:     b,
		Hearing:  contains("scheduled"),
		Status:   DetectState(b),
		BasePath: path,
	}
	if contains("FHV") {
		f.VehicleType = reg.FHV.String()
	}
	for _, r := range reg.All {
		if contains(r.Code) {
			f.Violations = append(f.Violations, r)
		}
	}
	f.Lines = splitLines(f.Body)
	if len(f.Lines) >= 1 {
		location := strings.SplitN(f.Lines[0], " ", 5)
		if len(location) == 5 {
			f.Location = location[4]
		}
	}
	if len(f.Lines) >= 2 {
		f.Description = f.Lines[1]
	}
	f.Tweets = tweetPattern.FindAllString(f.Body, -1)
	f.ServiceRequestID = findServiceRequestID(f.Lines)

	for _, filename := range files {
		switch filename {
		case "notes.txt", ".DS_Store":
			continue
		}
		ext := strings.ToLower(filepath.Ext(filename))
		switch ext {
		case ".mov", ".m4v":
			f.Videos = append(f.Videos, filename)
		case ".bmp", ".png", ".jpg", ".jpeg", ".tif", ".gif":
			f.Photos = append(f.Photos, filename)
		default:
			f.Files = append(f.Files, filename)
		}
	}

	return f, nil
}

func findServiceRequestID(lines []string) string {
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "1-1-1") && len(line) == 14:
			return "C" + line
		case strings.HasPrefix(line, "C1-1-1") && len(line) == 15:
			return line
		case strings.HasPrefix(line, "Service Request #: C1-1-1") && len(line) == 34:
			return line[19:34]
		}
	}
	return ""
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
