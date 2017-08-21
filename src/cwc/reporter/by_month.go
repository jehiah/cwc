package reporter

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"strings"

	"cwc/db"

	"github.com/olekukonko/tablewriter"
)

type ByMonth struct {
	counts map[string]int
	months []string
	Scale
}

func NewByMonth(d db.DB, f []*db.FullComplaint) (Reporter, error) {
	r := &ByMonth{
		counts: make(map[string]int),
	}
	for _, c := range f {
		month := c.Time.Format("200601")
		r.counts[month] += 1
		r.Scale.Update(r.counts[month])
	}
	for m, _ := range r.counts {
		r.months = append(r.months, m)
	}
	sort.Strings(r.months)
	return r, nil
}

func (r ByMonth) HTML() template.HTML {
	return ""
}

func (r ByMonth) Text() string {
	w := &bytes.Buffer{}

	// io.WriteString(w, "TLC Complaints by month\n")
	table := tablewriter.NewWriter(w)
	table.SetBorder(false)
	table.SetHeader([]string{"Month", "Complaints", r.Scale.String()})
	for _, month := range r.months {
		n := r.counts[month]
		table.Append([]string{month, fmt.Sprintf("%d", n), strings.Repeat("∎", n/r.Scale.Scale)})
	}
	table.Render()
	fmt.Fprint(w, "\n")
	return w.String()
}
