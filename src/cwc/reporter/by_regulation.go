package reporter

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"

	"cwc/db"
	"cwc/reg"
)

type regMatch struct {
	Key     string
	Count   int
	Percent float64
	Code    string
}

type byCount []regMatch

func (a byCount) Len() int           { return len(a) }
func (a byCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byCount) Less(i, j int) bool { return a[i].Count < a[j].Count }

type ByRegulation struct {
	Total   int
	MaxDesc int
	Matches []regMatch
}

func NewByRegulation(d db.DB, f []*db.FullComplaint) (Reporter, error) {
	r := &ByRegulation{
		Total: len(f),
	}

	short := func(r reg.Reg) string {
		if r.Short != "" {
			return r.Short
		}
		return r.Description
	}

	finder := func(needle string) []*db.FullComplaint {
		var out []*db.FullComplaint
		for _, c := range f {
			if c.Contains(needle) {
				out = append(out, c)
			}
		}
		return out
	}

	for _, reg := range reg.All {
		complaints := finder(reg.Code)
		count := len(complaints)
		if count == 0 {
			continue
		}
		m := regMatch{short(reg), count, (float64(count) / float64(r.Total)) * 100, reg.Code}
		if len(m.Key) > r.MaxDesc {
			r.MaxDesc = len(m.Key)
		}
		r.Matches = append(r.Matches, m)
	}
	sort.Sort(sort.Reverse(byCount(r.Matches)))
	return r, nil

}

var byRegulationHTML string = `
<div class="col-xs-6">
<h2>Violations Cited</h2>
<table class="table table-condensed">
<tbody>
	{{range .Matches}}
	<tr>
		<td class="text-right">{{.Key}}</td>
		<td class="number">{{.Count}}</td>
		<td class="number"><small>{{printf "%0.2f" .Percent}}%</small></td>
		<td><a href="./?q={{.Code}}">{{.Code}}</a></td>
	</tr>
	{{end}}
</tbody>
</table>
</div>
`
var byRegulationTemplate *template.Template = template.Must(template.New("foo").Parse(byRegulationHTML))

func (r ByRegulation) HTML() template.HTML {
	return GetTemplateString(byRegulationTemplate, r)
}

func (r ByRegulation) Text() string {
	w := &bytes.Buffer{}

	fmt.Fprintf(w, "Violations Cited:\n")
	size := fmt.Sprintf("%d", r.MaxDesc)
	for _, m := range r.Matches {
		percent := fmt.Sprintf("(%0.1f%%)", m.Percent)
		fmt.Fprintf(w, "%"+size+"s %3d %-8s %s\n", m.Key, m.Count, percent, m.Code)
	}
	fmt.Fprint(w, "\n")
	return w.String()
}
