package reporter

import (
	"net/http"
	"net/url"
	"net"
	"log"
	"os/exec"
	"time"
	
	"cwc/db"
)

func ReportServer(d db.DB) error {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, err := JSON(d)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write(body)
	})

	http.Handle("/", http.FileServer(http.Dir("html")))
	
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}
	go func() {
		time.Sleep(50*time.Millisecond)
		u := &url.URL{Scheme:"http", Host: listener.Addr().String()}
		log.Printf("Running cwc report server at %q", u)
		err = exec.Command("/usr/bin/open", u.String()).Run()
		if err != nil {
			log.Fatalf("%s", err)
			return
		}
	}()
	log.Fatal(http.Serve(listener, nil))
	return nil
}
