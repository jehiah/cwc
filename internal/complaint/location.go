package complaint

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
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
	if len(f.PhotoDetails) != len(f.Photos) {
		log.Printf("HasGPSInfo: photos not loaded %#v", f.Complaint)
	}
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
	if len(f.PhotoDetails) != len(f.Photos) {
		log.Printf("HasGPSInfo: photos not loaded %#v", f.Complaint)
	}
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

var geoclientCache map[string]LL

func geoclient(params *url.Values) LL {
	cacheKey := params.Get("borough") + params.Get("crossStreetOne") + params.Get("crossStreetTwo")
	if ll, ok := geoclientCache[cacheKey]; ok {
		return ll
	}
	// https://api.nyc.gov/geoclient/v1/doc/
	// /v1/intersection.json?crossStreetOne=broadway&crossStreetTwo=w 99 st&borough=manhattan&app_id=abc123&app_key=def456
	// curl -v  -X GET "https://api.cityofnewyork.us/geoclient/v1/intersection.json?app_id=...&app_key=...&crossStreetOne=West+25rd+St&crossStreetTwo=8th+Ave&borough=Manhattan"
	url := "https://api.nyc.gov/geo/geoclient/v1/intersection.json?" + params.Encode()

	log.Printf("GET %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("%s", err)
		return LL{}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("%s", err)
		return LL{}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%s", err)
		return LL{}
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
		return LL{}
	}
	ll := LL{
		Lat:  data.Intersection.Lat,
		Long: data.Intersection.Long,
	}
	geoclientCache[cacheKey] = ll
	return ll
}

// random produces a float [min, max]
// from https://groups.google.com/forum/#!topic/Golang-Nuts/_M-8hRpQs84
func random(min, max float64) float64 {
	if min > max {
		min, max = max, min
	}
	return rand.Float64()*(max-min) + min
}

// randomBetween picks a point on a line between a and bs
func randomBetween(a, b LL) LL {
	// for my mental sanity; pick the point on the left for the NYC area
	if b.Long < a.Long {
		a, b = b, a
	}
	diff := func(f1, f2 float64) float64 {
		return f2 - f1
	}
	midpoint := rand.Float64()
	ll := LL{
		Lat:  a.Lat + (diff(a.Lat, b.Lat) * midpoint),
		Long: a.Long + (diff(a.Long, b.Long) * midpoint),
	}
	// log.Printf("Lat: start %0.4f + (diff:%0.5f) * midpoint:%0.5f = %0.5f => %0.4f", a.Lat, diff(a.Lat, b.Lat), midpoint, (diff(a.Lat, b.Lat) * midpoint), ll.Lat)
	// log.Printf("Long: start %0.4f + (diff:%0.5f) * midpoint:%0.5f = %0.5f => %0.4f", a.Long, diff(a.Long, b.Long), midpoint, (diff(a.Long, b.Long) * midpoint), ll.Long)
	return ll
}

// GeoClientLookup does a lookup against the NYC 'geoclient'
// using cross streets to get a lat/long
// see: https://dev-mgmt.cityofnewyork.us/docs/geoclient/v1
func (f FullComplaint) GeoClientLookup() LL {
	street, crossStreet, betweenCrossStreet := ParseStreetCrossStreet(f.Location)

	AppID := os.Getenv("GEOCLIENT_APP_ID")
	AppKey := os.Getenv("GEOCLIENT_APP_KEY")
	if AppID == "" || AppKey == "" {
		return LL{}
	}

	params := &url.Values{
		"app_id":         []string{AppID},
		"app_key":        []string{AppKey},
		"borough":        []string{"Manhattan"},
		"crossStreetOne": []string{street},
		"crossStreetTwo": []string{crossStreet},
	}
	ll := geoclient(params)
	if ll.Lat == 0 {
		return ll
	}
	var tween bool
	if betweenCrossStreet != "" {
		params.Set("crossStreetTwo", betweenCrossStreet)
		ll2 := geoclient(params)

		if ll2.Lat != 0 {
			ll = randomBetween(ll, ll2)
			tween = true
		}
	}
	if !tween && ll.Lat != 0 {
		// nudge points slightly for randomization
		drift := 0.00096
		drift = 0.00032
		ll.Lat += random(-1*drift, drift)
		ll.Long += random(-1*drift, drift)
	}

	log.Printf("[%s] got %#v for %s %s %s", f.Complaint, ll, street, crossStreet, betweenCrossStreet)
	return ll
}

func ParseStreetCrossStreet(loc string) (s1, s2, s3 string) {
	loc = strings.Replace(loc, " Between ", " between ", -1)
	switch {
	case strings.Contains(loc, " between "):
		c := strings.SplitN(loc, " between ", 2)
		s1 = c[0]
		c = strings.SplitN(c[1], " and ", 2)
		s2 = c[0]
		s3 = c[1]
	case strings.Contains(loc, " and "):
		c := strings.SplitN(loc, " and ", 2)
		s1 = c[0]
		s2 = c[1]
	default:
		return loc, "", ""
	}
	return
}
