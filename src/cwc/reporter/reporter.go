package reporter

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"

	"cwc/db"
	"cwc/reg"
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

func JSON(d db.DB) ([]byte, error) {
	type record struct {
		Timestamp   int64     `json:"timestamp"`
		License     string    `json:"license_plate"`
		VehicleType string    `json:"vehicle_type"`
		Violations  []reg.Reg `json:"violations"`
		// Tweets []string `json:"tweets,omitempty"`
	}

	complaints, err := d.All()
	if err != nil {
		return nil, err
	}
	var data []*record
	for _, c := range complaints {
		record := &record{
			Timestamp:   c.Time().Unix(),
			License:     c.License(),
			VehicleType: "TAXI",
		}
		if ok, _ := d.ComplaintContains(c, " FHV "); ok {
			record.VehicleType = "FHV"
		}
		for _, r := range reg.All {
			if ok, _ := d.ComplaintContains(c, r.Code); ok {
				record.Violations = append(record.Violations, r)
			}
		}
		if len(record.Violations) == 0 {
			log.Printf("Warning: %s has no violations", c)
			continue
		}
		data = append(data, record)
	}
	body, err := json.Marshal(data)
	return body, err
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
