package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/jehiah/cwc/db"
	"github.com/jehiah/cwc/internal/complaint"
	"github.com/jehiah/cwc/internal/reg"
	"github.com/jehiah/cwc/internal/reporter"
	"github.com/jehiah/cwc/templates"
)

type Server struct {
	DB db.ReadWrite
	*template.Template
	*http.ServeMux
	listener net.Listener
	ReadOnly bool
	BasePath string

	medallionsMutex sync.Mutex
	medallions      []AuthorizedMedallion
}

func ComplaintClass(s complaint.State) string {
	switch s {
	case complaint.ClosedPenalty, complaint.ClosedInspection, complaint.NoticeOfDecision:
		return "success"
	case complaint.HearingScheduled:
		return "warning"
	case complaint.Fined:
		return "info"
	case complaint.ClosedUnableToID, complaint.Invalid, complaint.Expired:
		return "active"
	}
	return ""
}

func PhotoClass(p *complaint.Photo) string {
	switch p.Submitted {
	case true:
		return "panel-primary"
	case false:
		return "panel-info"
	}
	panic("here")
}

func New(d db.ReadWrite, templatePath, basePath string, readOnly bool) *Server {
	t := template.New("").Funcs(template.FuncMap{"ComplaintClass": ComplaintClass, "PhotoClass": PhotoClass})
	var err error
	if templatePath == "" {
		t, err = t.ParseFS(templates.FS, "*.html")
	} else {
		t, err = t.ParseGlob(filepath.Join(templatePath, "*.html"))
	}

	if err != nil {
		log.Fatalf("%s", err)
	}
	s := &Server{
		DB:       d,
		Template: t,
		ServeMux: http.NewServeMux(),
		ReadOnly: readOnly,
		BasePath: basePath,
	}
	s.ServeMux.HandleFunc(basePath+"reg", s.Regulations)
	s.ServeMux.HandleFunc(basePath+"complaint/", s.Complaint)
	s.ServeMux.HandleFunc("/", s.Complaints)
	s.ServeMux.HandleFunc(basePath+"data/report", s.DataReport)
	s.ServeMux.HandleFunc(basePath+"report", s.Report)
	s.ServeMux.HandleFunc(basePath+"report/taxi", s.TaxiReport)
	s.ServeMux.HandleFunc(basePath+"address.json", s.Address)
	return s
}

type AuthorizedMedallion struct {
	LicenseNumber string `json:"license_number"`
}

func (s *Server) AuthorizedMedallions() ([]AuthorizedMedallion, error) {
	s.medallionsMutex.Lock()
	defer s.medallionsMutex.Unlock()

	if len(s.medallions) > 1 {
		return s.medallions, nil
	}

	// get a list of authorized medallions - data for yesterday should always be available
	yesterday := time.Now().Add(time.Hour * -24).Format("2006-01-02")
	log.Printf("retrieving https://data.cityofnewyork.us/api/resource/rhe8-mgbb.json for %q", yesterday)
	resp, err := http.Get("https://data.cityofnewyork.us/api/resource/rhe8-mgbb.json?$select=license_number&$limit=100000&$where=last_updated_date+=+%27" + yesterday + "%27")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&s.medallions)
	if err != nil {
		return nil, err
	}
	return s.medallions, nil
}

func (s *Server) TaxiReport(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	query := r.Form.Get("q")
	var complaints []complaint.Complaint
	var err error
	if query == "" {
		complaints, err = s.DB.Index()
	} else {
		complaints, err = s.DB.Find(query)
	}

	if err != nil {
		log.Printf("%s", err)
		s.Error(w, err)
		return
	}
	seen := make(map[string]int)

	for _, c := range complaints {
		f, err := s.DB.FullComplaint(c)
		if err != nil {
			log.Printf("error parsing %s, %s", c, err)
			continue
		}
		if f.VehicleType != reg.Taxi.String() {
			continue
		}
		seen[f.License] += 1
	}

	medallions, err := s.AuthorizedMedallions()
	if err != nil {
		log.Printf("%s", err)
		s.Error(w, err)
		return
	}

	type Record struct {
		License string
		Count   int
	}
	all := make([]Record, 0, len(medallions))
	for _, license := range medallions {
		all = append(all, Record{License: license.LicenseNumber, Count: seen[license.LicenseNumber]})
	}
	sort.Slice(all, func(i, j int) bool { return strings.Compare(all[i].License, all[j].License) == -1 })

	err = s.Template.ExecuteTemplate(w, "report_taxi.html", all)
	if err != nil {
		log.Printf("%s", err)
		s.Error(w, err)
	}
}

func (s *Server) Report(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Query    string
		Page     string
		BasePath string
		Reports  []template.HTML
	}

	reports, err := reporter.RunHTML(s.DB)
	if err != nil {
		log.Printf("%s", err)
		s.Error(w, err)
		return
	}

	err = s.Template.ExecuteTemplate(w, "report.html", payload{Page: "Report", BasePath: s.BasePath, Reports: reports})
	if err != nil {
		log.Printf("%s", err)
		s.Error(w, err)
	}
}

// for /data/report
func (s *Server) DataReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var b bytes.Buffer
	err := reporter.JSON(&b, s.DB)
	if err != nil {
		log.Printf("%s", err)
		s.Error(w, err)
		return
	}
	w.Write(b.Bytes())
}

func (s *Server) Regulations(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Regulations []reg.Reg
		Query       string
		Page        string
		BasePath    string
	}
	p := payload{Regulations: reg.All, Page: "Regulations", BasePath: s.BasePath}
	err := s.Template.ExecuteTemplate(w, "reg.html", p)
	if err != nil {
		log.Printf("%s", err)
		s.Error(w, err)
	}
}

func (s *Server) OpenInBrowser() error {
	u := &url.URL{Scheme: "http", Host: s.listener.Addr().String(), Path: s.BasePath}
	err := exec.Command("/usr/bin/open", u.String()).Run()
	return err
}

func (s *Server) Serve(addr string, logRequests bool) error {
	if addr == "" {
		addr = ":53000"
	}
	var err error
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("Running cwc server at %s", s.listener.Addr())

	var h http.Handler = s
	if logRequests {
		h = handlers.LoggingHandler(os.Stdout, h)
	}

	err = http.Serve(s.listener, h)
	return err
}

func (s *Server) Complaints(w http.ResponseWriter, r *http.Request) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	if r.URL.Path != s.BasePath {
		http.NotFound(w, r)
		return
	}

	r.ParseForm()
	type payload struct {
		FullComplaints  []*complaint.FullComplaint
		PendingHearings []*complaint.FullComplaint
		Query           string
		Page            string
		BasePath        string
	}
	p := payload{
		Query:    r.Form.Get("q"),
		Page:     "Complaints",
		BasePath: s.BasePath,
	}
	var complaints []complaint.Complaint
	var err error
	if p.Query == "" || strings.HasPrefix(p.Query, "status:") || p.Query == "no_location" || p.Query == "no_discernable_location" {
		complaints, err = s.DB.Index()
	} else {
		complaints, err = s.DB.Find(p.Query)
	}
	if err != nil {
		log.Printf("%s", err)
		s.Error(w, err)
		return
	}
	for _, c := range complaints {
		f, err := s.DB.FullComplaint(c)
		if err != nil {
			log.Printf("error parsing %s, %s", c, err)
			continue
		}
		if strings.HasPrefix(p.Query, "status:") && f.Status.String() != p.Query[len("status:"):] {
			continue
		}
		if p.Query == "no_location" {
			if f.HasGPSInfo() {
				continue
			}
		}
		if p.Query == "no_discernable_location" {
			if f.HasGPSInfo() {
				continue
			}
			ll := f.GeoClientLookup()
			f.Lat, f.Long = ll.Lat, ll.Long
			if f.HasGPSInfo() {
				continue
			}
		}
		p.FullComplaints = append(p.FullComplaints, f)
	}

	if p.Query == "" {
		nyc, _ := time.LoadLocation("America/New_York")
		dayStart := time.Now().In(nyc).Truncate(time.Hour * 24)
		for _, f := range p.FullComplaints {
			if !f.Hearing.IsZero() && f.Hearing.After(dayStart) {
				p.PendingHearings = append(p.PendingHearings, f)
			}
		}
		sort.Sort(complaint.FullComplaintsByHearing(p.PendingHearings))
	}

	err = s.Template.ExecuteTemplate(w, "complaints.html", p)
	if err != nil {
		s.Error(w, err)
		log.Printf("%s", err)
	}
}

func (s *Server) Error(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	type payload struct {
		Error    string
		Query    string
		Page     string
		BasePath string
	}

	p := payload{Error: err.Error(), BasePath: s.BasePath}
	if s.ReadOnly {
		p.Error = "An Error Occurred"
	}
	err = s.Template.ExecuteTemplate(w, "error.html", p)
	if err != nil {
		log.Printf("error rendering %s", err)
	}
}

// Complaint handles
// /complaint/yyyymmdd_HHMM_LICENSE
// /complaint/yyyymmdd_HHMM_LICENSE.json
// /complaint/yyyymmdd_HHMM_LICENSE/iamge.jpg
func (s *Server) Complaint(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len(s.BasePath):]
	patterns := strings.SplitN(path, "/", 3)
	var isAPI bool
	if strings.HasSuffix(patterns[1], ".json") {
		isAPI = true
		patterns[1] = patterns[1][:len(patterns[1])-5]
	}
	c := complaint.Complaint(patterns[1])
	if patterns[1] == "latest" {
		var err error
		c, err = s.DB.Latest()
		if err != nil {
			http.Error(w, "Error getting latest", 500)
			return
		}
	}

	if ok, err := s.DB.Exists(c); err != nil {
		s.Error(w, err)
		return
	} else if !ok {
		http.NotFound(w, r)
		return
	}

	if len(patterns) == 3 {
		if r.Method != "GET" {
			http.Error(w, "Method not supported", 405)
			return
		}

		file := patterns[2]
		if file == "map" {
			s.Map(w, r, c)
			return
		}
		switch strings.ToLower(filepath.Ext(file)) {
		case ".png", ".jpg", ".jpeg", ".heic":
			s.Image(w, r, c, file)
		default:
			s.Download(w, r, c, file)
		}
		return
	}

	// handle POST
	if r.Method == "POST" {
		if s.ReadOnly {
			http.Error(w, "Method not supported", 405)
			return
		}
		r.ParseForm()
		txt := strings.TrimSpace(r.Form.Get("append_text"))
		if txt == "" {
			s.Error(w, fmt.Errorf("MISSING_ARG_APPEND_TEXT"))
			return
		}
		if !strings.HasPrefix(txt, "[ll") {
			txt = fmt.Sprintf("[note:%s] %s", time.Now().Format("2006/01/02 15:04"), txt)
		}
		s.DB.Append(c, fmt.Sprintf("\n%s\n", txt))
		http.Redirect(w, r, (&url.URL{Path: r.URL.Path}).String(), 302)
		return
	}

	f, err := s.DB.FullComplaint(c)
	if err != nil {
		s.Error(w, err)
		log.Printf("%s", err)
		return
	}

	f.PhotoDetails, err = db.LoadPhotos(s.DB, f)
	if err != nil {
		s.Error(w, err)
		log.Printf("%s", err)
		return
	}
	if isAPI {
		w.Header().Set("Content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		b, err := json.Marshal(JsonAPI(s.DB, f))
		if err != nil {
			http.Error(w, "JSON ERROR", 500)
			log.Printf("%s", err)
		} else {
			w.Write(b)
		}
		return
	}

	type payload struct {
		FullComplaint *complaint.FullComplaint
		Query         string
		Page          string
		ReadOnly      bool
		BasePath      string
	}
	p := payload{
		FullComplaint: f,
		Page:          "Complaints",
		ReadOnly:      s.ReadOnly,
		BasePath:      s.BasePath,
	}
	err = s.Template.ExecuteTemplate(w, "complaint.html", p)
	if err != nil {
		s.Error(w, err)
		log.Printf("%s", err)
	}
}

func (s *Server) Address(w http.ResponseWriter, r *http.Request) {
	if s.ReadOnly {
		http.Error(w, "not found", 404)
		return
	}
	type wrapper struct {
		Address struct {
			Email        string
			FirstName    string
			LastName     string
			PhoneNumber  string
			Borough      string
			AddressLine1 string
			AddressLine2 string
			City         string
			State        string
			ZipCode      string
		}
	}
	var o wrapper
	addrFile := s.DB.FullPath(complaint.Complaint("address.json"))
	f, err := os.Open(addrFile)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "not found", 404)
			return
		}
		http.Error(w, "error", 500)
		log.Printf("error %#f", err)
		return
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&o.Address)
	if err != nil {
		log.Printf("error %#f", err)
		http.Error(w, "error", 500)
		return
	}
	json.NewEncoder(w).Encode(o)
}

func JsonAPI(d db.ReadOnly, f *complaint.FullComplaint) interface{} {
	type wrapper struct {
		Complaint          *complaint.FullComplaint
		DateTimeOfIncident string
		DateTimeOfIncidentISO string
		Street             string
		CrossStreet        string
	}
	o := wrapper{
		Complaint: f,
		// TODO: move to JS formatting
		DateTimeOfIncident: f.Time.Format("1/2/2006 3:04 PM"),
		DateTimeOfIncidentISO: f.Time.Format("2006-01-02T15:04:05.000Z"), // 2023-05-26T21:08:00.0000000Z
	}
	o.Street, o.CrossStreet, _ = complaint.ParseStreetCrossStreet(o.Complaint.Location)
	return o
}
