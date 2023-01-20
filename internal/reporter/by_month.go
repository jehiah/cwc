package reporter

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/jehiah/cwc/db"
	"github.com/jehiah/cwc/internal/complaint"
	"github.com/olekukonko/tablewriter"
)

type ByMonth struct {
	Counts map[string]int
	Months []string
	Scale
}

func NewByMonth(d db.ReadOnly, f []*complaint.FullComplaint) (Reporter, error) {
	r := &ByMonth{
		Counts: make(map[string]int),
	}
	for _, c := range f {
		month := c.Time.Format("200601")
		r.Counts[month] += 1
		r.Scale.Update(r.Counts[month])
	}
	for m, _ := range r.Counts {
		r.Months = append(r.Months, m)
	}
	sort.Strings(r.Months)
	return r, nil
}

var byMonthHTML string = `
<div class="col-xs-6">
<h2>Complaints Per Month</h2>
<table class="table table-condensed">
<tbody>
	{{range .}}
	<tr>
		<td class="text-right">{{.Month}}</td>
		<td class="number">{{.N}}</td>
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
var byMonthTemplate *template.Template = template.Must(template.New("foo").Parse(byMonthHTML))

func (r ByMonth) HTML() template.HTML {
	type row struct {
		Month      string
		N          int
		Percent    float32
		RelPercent float32
	}
	var rows []row
	var max, total int
	for _, n := range r.Counts {
		total += n
		if n > max {
			max = n
		}
	}

	for _, month := range r.Months {
		n := r.Counts[month]
		rows = append(rows, row{month, n, percent(n, total), percent(n, max)})
	}
	return GetTemplateString(byMonthTemplate, rows)
}

func (r ByMonth) Text() string {
	w := &bytes.Buffer{}

	// io.WriteString(w, "TLC Complaints by month\n")
	table := tablewriter.NewWriter(w)
	table.SetBorder(false)
	table.SetHeader([]string{"Month", "Complaints", r.Scale.String()})
	for _, month := range r.Months {
		n := r.Counts[month]
		table.Append([]string{month, fmt.Sprintf("%d", n), strings.Repeat("∎", n/r.Scale.Scale)})
	}
	table.Render()
	fmt.Fprint(w, "\n")
	return w.String()
}
