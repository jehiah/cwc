package db

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// findLatLongLine parses a line: [ll:$lat,$long]
func findLatLongLine(lines []string) (lat, long float64) {
	// iterate backwards start with last log line
	for i := len(lines); i > 0; i-- {
		line := lines[i-1]
		if !strings.HasPrefix(line, "[ll:") {
			continue
		}
		line = line[4 : len(line)-1]
		chunks := strings.SplitN(line, ",", 2)
		if len(chunks) != 2 {
			continue
		}
		var err error
		lat, err = strconv.ParseFloat(strings.TrimSpace(chunks[0]), 64)
		if err != nil {
			continue
		}
		long, err = strconv.ParseFloat(strings.TrimSpace(chunks[1]), 64)
		if err != nil {
			continue
		}
		return
	}
	return
}

func (f *FullComplaint) HasGPSInfo() bool {
	if f.Lat != 0 && f.Long != 0 {
		return true
	}
	f.ParsePhotos()
	for _, p := range f.PhotoDetails {
		if p.Lat != 0 && p.Long != 0 {
			return true
		}
	}
	return false
}

func (f *FullComplaint) GPSInfo() LL {
	if f.Lat != 0 && f.Long != 0 {
		return LL{f.Lat, f.Long}
	}
	f.ParsePhotos()
	for _, p := range f.PhotoDetails {
		if p.Lat != 0 && p.Long != 0 {
			return LL{p.Lat, p.Long}
		}
	}
	return LL{}
}

// LL can be used from templates where you cant return (lat, long) separately from a function
type LL struct {
	Lat, Long float64
}

// GeoClientLookup does a lookup against the NYC 'geoclient'
// using cross streets to get a lat/long
// see: https://dev-mgmt.cityofnewyork.us/docs/geoclient/v1
func (f FullComplaint) GeoClientLookup() (ll LL) {
	street, crossStreet := ParseStreetCrossStreet(f.Location)

	AppID := os.Getenv("GEOCLIENT_APP_ID")
	AppKey := os.Getenv("GEOCLIENT_APP_KEY")
	if AppID == "" || AppKey == "" {
		return
	}

	params := &url.Values{
		"app_id":         []string{AppID},
		"app_key":        []string{AppKey},
		"borough":        []string{"Manhattan"},
		"crossStreetOne": []string{street},
		"crossStreetTwo": []string{crossStreet},
	}
	// /v1/intersection.json?crossStreetOne=broadway&crossStreetTwo=w 99 st&borough=manhattan&app_id=abc123&app_key=def456
	// curl -v  -X GET "https://api.cityofnewyork.us/geoclient/v1/intersection.json?app_id=...&app_key=...&crossStreetOne=West+25rd+St&crossStreetTwo=8th+Ave&borough=Manhattan"
	url := "https://api.cityofnewyork.us/geoclient/v1/intersection.json?" + params.Encode()
	log.Printf("GET %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	resp.Body.Close()

	type respBody struct {
		Intersection struct {
			Lat     float64 `json:"latitude"`
			Long    float64 `json:"longitude"`
			Message string  `json:"message"`
		} `json:"intersection"`
	}
	// log.Printf("> %s", string(body))
	var data respBody
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	ll = LL{
		Lat:  data.Intersection.Lat,
		Long: data.Intersection.Long,
	}
	log.Printf("got %#v for %s %s", ll, street, crossStreet)
	return ll
}

func ParseStreetCrossStreet(loc string) (s1, s2 string) {
	loc = strings.Replace(loc, " Between ", " between ", -1)
	switch {
	case strings.Contains(loc, " between "):
		c := strings.SplitN(loc, " between ", 2)
		s1 = c[0]
		c = strings.SplitN(c[1], " and ", 2)
		s2 = c[0]
	case strings.Contains(loc, " and "):
		c := strings.SplitN(loc, " and ", 2)
		s1 = c[0]
		s2 = c[1]
	default:
		return loc, ""
	}
	return
}
