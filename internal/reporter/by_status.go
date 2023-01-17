package reporter

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/jehiah/cwc/internal/complaint"
	"github.com/jehiah/cwc/internal/db"
)

type ByStatus struct {
	Data map[complaint.State]int
	Scale
	Total int
}

func NewByStatus(d db.ReadOnly, f []*complaint.FullComplaint) (Reporter, error) {
	r := &ByStatus{
		Data:  make(map[complaint.State]int),
		Total: len(f),
	}

	for _, c := range f {
		r.Data[c.Status]++
		r.Scale.Update(r.Data[c.Status])
	}
	return r, nil
}

var byStatusHTML string = `
<div class="col-xs-6">
<h2>Complaint Status</h2>
<table class="table table-condensed">
<tbody>
	{{range .}}
	<tr>
		<td ><a href="./?q=status:{{.State}}">{{.State}}</a></td>
		<td class="number">{{.N}}</td>
		<td class="number"><small>{{printf "%0.2f" .Percent}}%</small></td>
		<td>
		<div class="progress">
		  <div class="progress-bar {{.Class}}" role="progressbar" aria-valuenow="{{printf "%0.2f" .Percent}}" aria-valuemin="0" aria-valuemax="100" style="width: {{printf "%0.2f" .RelPercent}}%;"></div>
		</div>
		</td>
	</tr>
	{{end}}
</tbody>
</table>
</div>
`
var byStatusTemplate *template.Template = template.Must(template.New("foo").Parse(byStatusHTML))

func (r ByStatus) HTML() template.HTML {
	type row struct {
		State      string
		Class      string
		N          int
		Percent    float32
		RelPercent float32
	}
	var rows []row
	var max, total int

	for _, state := range complaint.AllStates {
		n := r.Data[state]
		total += n
		if n > max {
			max = n
		}
	}

	for _, state := range complaint.AllStates {
		n := r.Data[state]
		if n == 0 {
			continue
		}
		class := ComplaintClass(state)
		if class != "" {
			class = "progress-bar-" + class
		}
		rows = append(rows, row{string(state), class, n, percent(n, total), percent(n, max)})
	}
	return GetTemplateString(byStatusTemplate, rows)
}

func ComplaintClass(s complaint.State) string {
	switch s {
	case complaint.ClosedPenalty, complaint.ClosedInspection, complaint.NoticeOfDecision:
		return "success"
	case complaint.HearingScheduled:
		return "warning"
	case complaint.Fined:
		return "info"
	case complaint.ClosedUnableToID, complaint.Invalid, complaint.Expired:
		return "active"
	}
	return ""
}

func (r ByStatus) Text() string {
	w := &bytes.Buffer{}

	for _, state := range complaint.AllStates {
		n := r.Data[state]
		if n == 0 {
			continue
		}
		fmt.Fprintf(w, "%24s [ %3d complaints] %s (%0.1f%%)\n", state, n, strings.Repeat("∎", n/r.Scale.Scale), percent(n, r.Total))
	}
	fmt.Fprintf(w, "%24s [ %3d complaints] %s", "Total:", r.Total, r.Scale.String())
	fmt.Fprint(w, "\n")
	return w.String()
}
