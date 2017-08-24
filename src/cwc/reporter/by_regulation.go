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
	Percent float32
	Codes   []string
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

	lookup := make(map[string]*regMatch)
	for _, c := range f {
		for _, v := range c.Violations {
			m, ok := lookup[short(v)]
			if !ok {
				m = &regMatch{Key: short(v)}
				lookup[short(v)] = m
			}
			m.Count += 1
			ok = false
			for _, code := range m.Codes {
				if code == v.Code {
					ok = true
				}
			}
			if !ok {
				m.Codes = append(m.Codes, v.Code)
			}
		}
	}
	for _, m := range lookup {
		m.Percent = percent(m.Count, r.Total)
		r.Matches = append(r.Matches, *m)
		if len(m.Key) > r.MaxDesc {
			r.MaxDesc = len(m.Key)
		}
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
		<td><small>{{range $i, $c := .Codes }}{{if $i}}, {{end}}<a href="./?q={{$c}}">{{$c}}</a>{{end}}</small></td>
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
		fmt.Fprintf(w, "%"+size+"s %3d %-8s ", m.Key, m.Count, percent)
		for i, c := range m.Codes {
			if i == 0 {
				fmt.Fprintf(w, "%s", c)
			} else {
				fmt.Fprintf(w, ", %s", c)
			}
		}
		fmt.Fprint(w, "\n")
	}
	fmt.Fprint(w, "\n")
	return w.String()
}
