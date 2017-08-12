package reporter

import (
	"fmt"
	"io"
	"strings"

	"cwc/db"
)

func ByStatus(d db.DB, w io.Writer) error {
	data := make(map[db.State]int)
	complaints, err := d.All()
	if err != nil {
		return err
	}
	var total int
	scale := &Scale{}
	for _, c := range complaints {
		f, err := d.FullComplaint(c)
		if err != nil {
			continue
		}
		total++
		data[f.Status]++
		scale.Update(data[f.Status])
	}

	io.WriteString(w, scale.String())

	for _, state := range db.AllStates {
		n := data[state]
		if n == 0 {
			continue
		}
		p := (float32(n) / float32(total)) * 100
		fmt.Fprintf(w, "%30s [ %3d complaints] %s (%0.1f%%)\n", state, n, strings.Repeat("∎", n/scale.Scale), p)
	}
	fmt.Fprint(w, "\n")
	return nil
}
