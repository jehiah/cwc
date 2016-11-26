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
}

func New(d db.DB, templatePath string) *Server {
	t, err := template.New("").Funcs(template.FuncMap{"ComplaintClass": func(c *db.FullComplaint) string {
		switch c.Status {
		case db.ClosedPenalty, db.ClosedInspection:
			return "success"
		case db.HearingScheduled:
			return "warning"
		case db.Fined:
			return "info"
		}
		return ""
	}}).ParseGlob(filepath.Join(templatePath, "*.html"))

	if err != nil {
		log.Fatalf("%s", err)
	}
	s := &Server{
		DB:       d,
		Template: t,
		ServeMux: http.NewServeMux(),
	}
	s.ServeMux.HandleFunc("/reg", s.Regulations)
	s.ServeMux.HandleFunc("/complaint/", s.Complaint)
	s.ServeMux.HandleFunc("/", s.Complaints)
	s.ServeMux.HandleFunc("/data/report", s.DataReport)
	s.ServeMux.HandleFunc("/report", s.Report)
	return s
}

func (s *Server) Report(w http.ResponseWriter, r *http.Request) {
	err := s.Template.ExecuteTemplate(w, "report.html", nil)
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
	}
	p := payload{Regulations: reg.All}
	err := s.Template.ExecuteTemplate(w, "reg.html", p)
	if err != nil {
		log.Printf("%s", err)
		s.Error(w, err)
	}
}

func (s *Server) OpenInBrowser() error {
	u := &url.URL{Scheme: "http", Host: s.listener.Addr().String()}
	err := exec.Command("/usr/bin/open", u.String()).Run()
	return err
}

func (s *Server) Serve(addr string) error {
	if addr == "" {
		addr = ":0"
	}
	var err error
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("Running cwc server at %s", s.listener.Addr())
	go func() {
		time.Sleep(200 * time.Millisecond)
		err := s.OpenInBrowser()
		if err != nil {
			log.Println(err)
		}
	}()
	err = http.Serve(s.listener, s)
	return err
}

func (s *Server) Complaints(w http.ResponseWriter, r *http.Request) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	r.ParseForm()
	type payload struct {
		FullComplaints []*db.FullComplaint
		Query          string
	}
	p := payload{
		Query: r.Form.Get("q"),
	}
	var complaints []db.Complaint
	var err error
	if p.Query == "" {
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
		p.FullComplaints = append(p.FullComplaints, f)
	}
	err = s.Template.ExecuteTemplate(w, "complaints.html", p)
	if err != nil {
		s.Error(w, err)
		log.Printf("%s", err)
	}
}

func (s *Server) Error(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	err = s.Template.ExecuteTemplate(w, "error.html", struct{ Error string }{err.Error()})
	if err != nil {
		log.Printf("error rendering %s", err)
	}
}

func (s *Server) Complaint(w http.ResponseWriter, r *http.Request) {
	patterns := strings.SplitN(r.URL.Path[1:], "/", 3)

	c := db.Complaint(patterns[1])
	f, err := s.DB.FullComplaint(c)
	if err != nil {
		s.Error(w, err)
		log.Printf("%s", err)
		return
	}

	if len(patterns) == 3 {
		file := patterns[2]
		var found bool
		for _, f := range f.Photos {
			if f == file {
				found = true
			}
		}
		for _, f := range f.Files {
			if f == file {
				found = true
			}
		}
		for _, f := range f.Videos {
			if f == file {
				found = true
			}
		}
		if found {
			path := s.DB.FullPath(c)
			staticServer := http.StripPrefix(fmt.Sprintf("/complaint/%s/", patterns[1]), http.FileServer(http.Dir(path)))
			staticServer.ServeHTTP(w, r)
			return
		}
		log.Printf("temp 404 %q", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	r.ParseForm()

	type payload struct {
		FullComplaint *db.FullComplaint
	}
	p := payload{
		FullComplaint: f,
	}
	err = s.Template.ExecuteTemplate(w, "complaint.html", p)
	if err != nil {
		s.Error(w, err)
		log.Printf("%s", err)
	}

}
