package reporter

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"

	"cwc/db"
)

type Reporter interface {
	HTML() template.HTML
	Text() string
}
type Generator func(db.DB, []*db.FullComplaint) (Reporter, error)

func Run(d db.DB, w io.Writer) error {
	var full []*db.FullComplaint
	complaints, err := d.All()
	if err != nil {
		return err
	}
	for _, c := range complaints {
		f, err := d.FullComplaint(c)
		if err != nil {
			continue
		}
		full = append(full, f)
	}

	for _, g := range []Generator{NewByMonth, NewByHour, NewPerDay, NewByStatus, NewByRegulation, NewByVehicle} {
		r, err := g(d, full)
		if err != nil {
			return err
		}
		io.WriteString(w, r.Text())
	}
	return nil
}

func RunHTML(d db.DB) ([]template.HTML, error) {
	var full []*db.FullComplaint
	complaints, err := d.All()
	if err != nil {
		return nil, err
	}
	for _, c := range complaints {
		f, err := d.FullComplaint(c)
		if err != nil {
			continue
		}
		full = append(full, f)
	}

	var o []template.HTML
	for _, g := range []Generator{NewByMonth, NewByHour, NewPerDay, NewByStatus, NewByRegulation, NewByVehicle} {
		r, err := g(d, full)
		if err != nil {
			return nil, err
		}
		h := r.HTML()
		if h == "" {
			continue
		}
		o = append(o, h)
	}
	return o, nil
}

func JSON(w io.Writer, d db.DB) error {
	complaints, err := d.All()
	if err != nil {
		return err
	}
	e := json.NewEncoder(w)
	e.SetEscapeHTML(false)

	for _, c := range complaints {
		fc, err := d.FullComplaint(c)
		if err != nil {
			return err
		}

		if len(fc.Violations) == 0 {
			log.Printf("Warning: %s has no violations", c)
			continue
		}

		fc.ParsePhotos()
		if !fc.HasGPSInfo() {
			ll := fc.GeoClientLookup()
			fc.Lat, fc.Long = ll.Lat, ll.Long
		}
		err = e.Encode(fc)
		if err != nil {
			return err
		}
	}
	return nil
}

func percent(n, total int) float32 {
	return (float32(n) / float32(total)) * 100
}
func GetTemplateString(t *template.Template, data interface{}) template.HTML {
	b := &bytes.Buffer{}
	err := t.Execute(b, data)
	if err != nil {
		panic(err)
	}
	return template.HTML(b.String())
}
