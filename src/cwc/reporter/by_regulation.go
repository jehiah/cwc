package reporter

import (
	"fmt"
	"io"
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

func ByRegulation(d db.DB, w io.Writer) error {
	allComplaints, err := d.All()
	if err != nil {
		return err
	}
	totalComplaints := float64(len(allComplaints))

	short := func(r reg.Reg) string {
		if r.Short != "" {
			return r.Short
		}
		return r.Description
	}

	var matches []regMatch
	maxDesc := 0
	for _, r := range reg.All {
		complaints, err := d.Find(r.Code)
		if err != nil {
			return err
		}
		count := len(complaints)
		if count == 0 {
			continue
		}
		m := regMatch{short(r), count, (float64(count) / totalComplaints) * 100, r.Code}
		if len(m.Key) > maxDesc {
			maxDesc = len(m.Key)
		}
		matches = append(matches, m)
	}
	sort.Sort(byCount(matches))

	fmt.Fprintf(w, "Violations Cited:\n")
	size := fmt.Sprintf("%d", maxDesc)
	for _, m := range matches {
		percent := fmt.Sprintf("(%0.1f%%)", m.Percent)
		fmt.Fprintf(w, "%"+size+"s %3d %-8s %s\n", m.Key, m.Count, percent, m.Code)
	}
	fmt.Fprint(w, "\n")
	return nil
}
