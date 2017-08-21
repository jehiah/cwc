package reporter

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"cwc/db"

	"github.com/olekukonko/tablewriter"
)

type ByHour struct {
	Scale
	Start, Stop, Total int
	Hours              [24]int
}

func NewByHour(d db.DB, f []*db.FullComplaint) (Reporter, error) {
	r := &ByHour{
		Total: len(f),
	}

	for _, c := range f {
		hour := c.Time.Hour()
		r.Hours[hour] += 1
		r.Scale.Update(r.Hours[hour])
	}

	for i, v := range r.Hours {
		if (r.Start == 0 || r.Start == i-1) && v == 0 {
			r.Start = i
		}
		if v != 0 {
			r.Stop = i
		}
	}
	return r, nil
}

func (r ByHour) HTML() template.HTML {
	return ""
}

func (r ByHour) Text() string {
	w := &bytes.Buffer{}

	table := tablewriter.NewWriter(w)
	table.SetBorder(false)
	table.SetHeader([]string{"Hour", "Complaints", strings.TrimSpace(r.Scale.String())})
	for h, v := range r.Hours[r.Start+1 : r.Stop] {
		t := time.Date(2010, 1, 1, h+r.Start, 0, 0, 0, time.UTC)
		table.Append([]string{t.Format("3pm"), fmt.Sprintf("%d", v), strings.Repeat("∎", v/r.Scale.Scale)})
	}
	table.SetFooter([]string{"Total:", fmt.Sprintf("%d", r.Total), ""})
	table.Render()
	fmt.Fprint(w, "\n")
	return w.String()
}
