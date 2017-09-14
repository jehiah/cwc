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

var byHourHTML string = `
<div class="col-xs-6">
<h2>Complaints Per Hour</h2>
<table class="table table-condensed">
<thead>
	<tr>
		<th class="number">Hour</th>
		<th class="number">Complaints</th>
		<th></th>
	</tr>
</thead>
<tbody>
	{{range .}}
	<tr>
		<td class="number">{{.T.Format "3pm"}}</td>
		<td class="number">{{.N}}</td>
		<td>
		<div class="progress">
		  <div class="progress-bar" role="progressbar" aria-valuenow="{{printf "%0.2f" .Percent}}" aria-valuemin="0" aria-valuemax="100" style="width: {{printf "%0.2f" .RelPercent}}%;">
		    {{printf "%0.2f" .Percent}}%
		  </div>
		</div>
		</td>
	</tr>
	{{end}}
</tbody>
</table>
</div>
`
var byHourTemplate *template.Template = template.Must(template.New("foo").Parse(byHourHTML))

func (r ByHour) HTML() template.HTML {
	type row struct {
		T          time.Time
		N          int
		Percent    float32
		RelPercent float32
	}
	var rows []row
	var max int
	for _, n := range r.Hours {
		if n > max {
			max = n
		}
	}

	for h, n := range r.Hours[r.Start+1 : r.Stop] {
		t := time.Date(2010, 1, 1, h+r.Start+1, 0, 0, 0, time.UTC)
		rows = append(rows, row{t, n, percent(n, r.Total), percent(n, max)})
	}

	return GetTemplateString(byHourTemplate, rows)
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
