package reporter

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"strings"

	"cwc/db"
)

type PerDay struct {
	frequency []int
	Scale
}

func NewPerDay(d db.DB, f []*db.FullComplaint) (Reporter, error) {
	r := &PerDay{}

	counts := make(map[string]int)

	var max int
	for _, c := range f {
		date := c.Time.Format("20060102")
		counts[date] += 1
		if counts[date] > max {
			max = counts[date]
		}
	}
	if max == 0 {
		return nil, nil
	}

	r.frequency = make([]int, max)
	for _, n := range counts {
		r.frequency[n-1] += 1
		r.Scale.Update(r.frequency[n-1])
	}
	return r, nil
}

func (r PerDay) HTML() template.HTML {
	return ""
}

func (r PerDay) Text() string {
	w := &bytes.Buffer{}

	io.WriteString(w, "Distribution of Complaints per day:\n")
	io.WriteString(w, r.Scale.String())
	for freq, n := range r.frequency {
		if n == 0 {
			continue
		}
		fmt.Fprintf(w, "%2d complaints/day [ %3d days] %s\n", freq+1, n, strings.Repeat("∎", n/r.Scale.Scale))
	}
	fmt.Fprint(w, "\n")
	return w.String()
}
