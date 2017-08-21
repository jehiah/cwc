package reporter

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"cwc/db"
)

type ByStatus struct {
	Data map[db.State]int
	Scale
	Total int
}

func NewByStatus(d db.DB, f []*db.FullComplaint) (Reporter, error) {
	r := &ByStatus{
		Data:  make(map[db.State]int),
		Total: len(f),
	}

	for _, c := range f {
		r.Data[c.Status]++
		r.Scale.Update(r.Data[c.Status])
	}
	return r, nil
}

func (r ByStatus) HTML() template.HTML {
	return ""
}

func (r ByStatus) Text() string {
	w := &bytes.Buffer{}

	for _, state := range db.AllStates {
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
