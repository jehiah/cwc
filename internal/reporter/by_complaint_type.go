package reporter

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/jehiah/cwc/internal/db"
)

type byViolation struct {
	Type    string
	Count   int
	Percent float32
}

type ByViolationType struct {
	Total   int
	Matches map[string]*byViolation
}

func NewByViolationType(d db.DB, f []*db.FullComplaint) (Reporter, error) {
	r := &ByViolationType{
		Total: len(f),
		Matches: map[string]*byViolation{
			"parking": &byViolation{Type: "parking"},
			"moving":  &byViolation{Type: "moving"},
			"other":   &byViolation{Type: "other"},
		},
	}

	for _, c := range f {
		var violationType string
		for _, v := range c.Violations {
			violationType = v.Type
			if v.Type == "moving" {
				break
			}
		}
		if violationType == "" {
			continue
		}
		r.Matches[violationType].Count += 1
	}
	for _, m := range r.Matches {
		m.Percent = percent(m.Count, r.Total)
	}
	return r, nil
}

var byViolationTypeHTML string = `
<div class="col-xs-6">
<h2>Violations Type</h2>
<table class="table table-condensed">
<tbody>
	{{range .Matches}}
	<tr>
		<td class="text-right">{{.Type}}</td>
		<td class="number">{{.Count}}</td>
		<td class="number"><small>{{printf "%0.2f" .Percent}}%</small></td>
	</tr>
	{{end}}
</tbody>
</table>
</div>
`
var byViolationTypeTemplate *template.Template = template.Must(template.New("foo").Parse(byViolationTypeHTML))

func (r ByViolationType) HTML() template.HTML {
	return GetTemplateString(byViolationTypeTemplate, r)
}

func (r ByViolationType) Text() string {
	w := &bytes.Buffer{}

	fmt.Fprintf(w, "Violation Type:\n")
	for _, m := range r.Matches {
		percent := fmt.Sprintf("(%0.1f%%)", m.Percent)
		fmt.Fprintf(w, "%s %d %s", m.Type, m.Count, percent)
		fmt.Fprint(w, "\n")
	}
	fmt.Fprint(w, "\n")
	return w.String()
}
