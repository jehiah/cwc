package reporter

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/jehiah/cwc/db"
	"github.com/jehiah/cwc/internal/complaint"
)

type PerDay struct {
	frequency []int
	Scale
}

func NewPerDay(d db.ReadOnly, f []*complaint.FullComplaint) (Reporter, error) {
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

var perDayHTML string = `
<div class="col-xs-6">
<h2>Daily Complaint Frequency</h2>
<table class="table table-condensed">
<tbody>
	{{range .}}
	<tr>
		<td class="text-right">{{.Frequency}} complaints/day</td>
		<td class="text-right">{{.N}} days</td>
		<td>
		<div class="progress">
		  <div class="progress-bar" role="progressbar" aria-valuenow="{{printf "%0.2f" .Percent}}" aria-valuemin="0" aria-valuemax="100" style="width: {{printf "%0.2f" .RelPercent}}%;"></div>
		</div>
		</td>
	</tr>
	{{end}}
</tbody>
</table>
</div>
`
var perDayTemplate *template.Template = template.Must(template.New("foo").Parse(perDayHTML))

func (r PerDay) HTML() template.HTML {
	type row struct {
		Frequency  int
		N          int
		Percent    float32
		RelPercent float32
	}
	var rows []row
	var max, total int
	for _, n := range r.frequency {
		total += n
		if n > max {
			max = n
		}
	}

	for freq, n := range r.frequency {
		if n == 0 {
			continue
		}
		rows = append(rows, row{freq + 1, n, percent(n, total), percent(n, max)})
	}
	return GetTemplateString(perDayTemplate, rows)
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
