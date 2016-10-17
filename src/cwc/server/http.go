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
	"time"

	"cwc/db"
	"cwc/reg"
	"cwc/reporter"
)

type Server struct {
	db.DB
	*template.Template

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
	return &Server{
		DB:       d,
		Template: t,
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

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.URL.Path {
	case "/":
		err = s.Complaints(w, r)
	case "/reg":
		type payload struct {
			Regulations []reg.Reg
		}
		p := payload{Regulations: reg.All}
		err = s.Template.ExecuteTemplate(w, "reg.html", p)
	case "/report":
		err = s.Template.ExecuteTemplate(w, "report.html", nil)
	case "/data/report":
		w.Header().Set("Content-Type", "application/json")
		var body []byte
		body, err = reporter.JSON(s.DB)
		if err != nil {
			break
		}
		w.Write(body)
	default:
		http.Error(w, fmt.Sprintf("path %q not found", r.URL.Path), 404)
		return
	}
	if err != nil {
		log.Printf("%s", err)
		http.Error(w, "Unknown Error", 500)
		return
	}
}

func (s *Server) Complaints(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	type payload struct {
		FullComplaints []*db.FullComplaint
		Query string
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
		return err
	}
	for _, c := range complaints {
		f, err := s.DB.FullComplaint(c)
		if err != nil {
			log.Printf("error parsing %s, %s", c, err)
			continue
		}
		p.FullComplaints = append(p.FullComplaints, f)
	}
	return s.Template.ExecuteTemplate(w, "complaints.html", p)
}
