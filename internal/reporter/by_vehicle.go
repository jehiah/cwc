package reporter

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/jehiah/cwc/internal/complaint"
	"github.com/jehiah/cwc/internal/db"
)

type ByVehicle struct {
	Repeats          map[string][]time.Time
	FHV, Taxi, Total int
}

func NewByVehicle(d db.ReadOnly, f []*complaint.FullComplaint) (Reporter, error) {
	r := &ByVehicle{
		Repeats: make(map[string][]time.Time),
		Total:   len(f),
	}
	for _, c := range f {
		r.Repeats[c.License] = append(r.Repeats[c.License], c.Time)
		if c.Contains(" FHV ") {
			r.FHV++
		} else {
			r.Taxi++
		}
	}
	return r, nil
}

func (r ByVehicle) HTML() template.HTML {
	return ""
}

func (r ByVehicle) Text() string {
	w := &bytes.Buffer{}

	var preamble bool

	var repeatCount int
	for _, times := range r.Repeats {
		if len(times) < 2 {
			continue
		}
		repeatCount++
	}

	if repeatCount > 0 {
		fmt.Fprintf(w, "License Plates w/ Multiple Reports: %d\n", repeatCount)
		preamble = true
	}

	for l, times := range r.Repeats {
		if len(times) < 2 {
			continue
		}
		var suffix []string
		for _, t := range times {
			suffix = append(suffix, t.Format("2006-01-02"))
		}
		fmt.Fprintf(w, "%-7s seen %d times (%s)\n", l, len(times), strings.Join(suffix, ", "))
	}
	if preamble {
		fmt.Fprint(w, "\n")
	}

	totalLicenseCount := len(r.Repeats)
	fmt.Fprintf(w, "Number of Unique License Plates: %d (of %d reports)\n", totalLicenseCount, r.Total)
	fmt.Fprintf(w, "Taxi: %d (%0.1f%%) FHV: %d (%0.1f%%)\n", r.Taxi, percent(r.Taxi, totalLicenseCount), r.FHV, percent(r.FHV, totalLicenseCount))
	fmt.Fprint(w, "\n")

	return w.String()
}
