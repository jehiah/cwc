package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/jehiah/cwc/internal/db"
)

func (s *Server) Map(w http.ResponseWriter, r *http.Request, c db.Complaint) {
	f, err := s.DB.FullComplaint(c)
	if err != nil {
		http.Error(w, "UNKNOWN_ERROR", 500)
		log.Printf("%s", err)
		return
	}

	// https://api.mapbox.com/styles/v1/mapbox/streets-v8/static/-122.4241,37.78,14.25,-10,0/600x600?access_token=
	// env.Get("MAPBOX_TOKEN")
	accessToken := os.Getenv("MAPBOX_TOKEN")
	if accessToken == "" {
		http.Error(w, "Mapbox not configured", 404)
		return
	}
	ll := f.GPSInfo()

	r.ParseForm()
	size := r.Form.Get("s")
	if size == "" {
		size = "600x600"
	}
	zoom := r.Form.Get("z")
	if zoom == "" {
		zoom = "15"
	}

	rotation := 28 // the manhattan street grid offset
	tile := fmt.Sprintf("%0.4f,%0.4f,%s,%d,0", ll.Long, ll.Lat, zoom, rotation)
	params := url.Values{"access_token": {accessToken}}
	url := &url.URL{
		Scheme:   "https",
		Host:     "api.mapbox.com",
		Path:     fmt.Sprintf("/styles/v1/mapbox/streets-v8/static/%s/%s@2x", tile, size),
		RawQuery: params.Encode(),
	}
	http.Redirect(w, r, url.String(), 302)
}
