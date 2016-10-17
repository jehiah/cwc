package db

import (
	"bytes"
	"cwc/reg"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

type FullComplaint struct {
	Timestamp   int64     `json:"timestamp"`
	Time        time.Time `json:"-"`
	License     string    `json:"license_plate"`
	VehicleType string    `json:"vehicle_type"`
	Location    string    `json:"location"`
	Description string    `json:"description"`

	Status           string `json:"status"`
	ServiceRequestID string `json:"311_service_request_id"`
	TLCID            string `json:"tlc_id,omitempty"`

	Body       string    `json:"body"`
	Violations []reg.Reg `json:"violations"`
	Tweets     []string  `json:"tweets,omitempty"`
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
	return ParseComplaint(c, body)
}

var tweetPattern = regexp.MustCompile(`https?://[^\s]+`)

func ParseComplaint(c Complaint, body []byte) (*FullComplaint, error) {
	f := &FullComplaint{
		Timestamp:   c.Time().Unix(),
		Time:        c.Time(),
		License:     c.License(),
		VehicleType: reg.Taxi.String(),

		Body: string(body),
	}
	contains := func(pattern string) bool {
		return bytes.Contains(body, []byte(pattern))
	}
	if contains("FHV") {
		f.VehicleType = reg.FHV.String()
	}
	for _, r := range reg.All {
		if contains(r.Code) {
			f.Violations = append(f.Violations, r)
		}
	}
	lines := splitLines(f.Body)
	if len(lines) >= 1 {
		location := strings.SplitN(lines[0], " ", 5)
		if len(location) == 5 {
			f.Location = location[4]
		}
	}
	if len(lines) >= 2 {
		f.Description = lines[1]
	}
	f.Tweets = tweetPattern.FindAllString(f.Body, -1)
	f.ServiceRequestID = findServiceRequestID(lines)
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
