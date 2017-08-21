package server

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"cwc/db"
	"cwc/reg"
	"cwc/reporter"
)

type Server struct {
	db.DB
	*template.Template
	*http.ServeMux
	listener net.Listener
	ReadOnly bool
	BasePath string
}

func ComplaintClass(s db.State) string {
	switch s {
	case db.ClosedPenalty, db.ClosedInspection, db.NoticeOfDecision:
		return "success"
	case db.HearingScheduled:
		return "warning"
	case db.Fined:
		return "info"
	case db.ClosedUnableToID, db.Invalid, db.Expired:
		return "active"
	}
	return ""
}

func PhotoClass(p *db.Photo) string {
	switch p.Submitted {
	case true:
		return "panel-primary"
	case false:
		return "panel-info"
	}
	panic("here")
}

func New(d db.DB, templatePath, basePath string, readOnly bool) *Server {
	t, err := template.New("").Funcs(template.FuncMap{"ComplaintClass": ComplaintClass, "PhotoClass": PhotoClass}).ParseGlob(filepath.Join(templatePath, "*.html"))

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
	return s
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
	body, err := reporter.JSON(s.DB)
	if err != nil {
		log.Printf("%s", err)
		s.Error(w, err)
		return
	}
	w.Write(body)
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

func (s *Server) Serve(addr string) error {
	if addr == "" {
		addr = ":53000"
	}
	var err error
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("Running cwc server at %s", s.listener.Addr())
	err = http.Serve(s.listener, s)
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
		FullComplaints  []*db.FullComplaint
		PendingHearings []*db.FullComplaint
		Query           string
		Page            string
		BasePath        string
	}
	p := payload{
		Query:    r.Form.Get("q"),
		Page:     "Complaints",
		BasePath: s.BasePath,
	}
	var complaints []db.Complaint
	var err error
	if p.Query == "" || strings.HasPrefix(p.Query, "status:") {
		complaints, err = s.DB.All()
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
		sort.Sort(db.FullComplaintsByHearing(p.PendingHearings))
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

func (s *Server) Complaint(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len(s.BasePath):]
	patterns := strings.SplitN(path, "/", 3)
	c := db.Complaint(patterns[1])

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
		case ".png", ".jpg", ".jpeg":
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
		s.DB.Append(c, fmt.Sprintf("\n[note:%s] %s\n", time.Now().Format("2006/01/02 15:04"), txt))
		http.Redirect(w, r, (&url.URL{Path: r.URL.Path}).String(), 302)
		return
	}

	f, err := s.DB.FullComplaint(c)
	if err != nil {
		s.Error(w, err)
		log.Printf("%s", err)
		return
	}

	f.ParsePhotos()

	type payload struct {
		FullComplaint *db.FullComplaint
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
