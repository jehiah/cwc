package reporter

import (
	"fmt"
	"io"
	"strings"

	"cwc/db"
)

func PerDay(d db.DB, w io.Writer) error {

	counts := make(map[string]int)
	complaints, err := d.All()
	if err != nil {
		return err
	}
	var max int
	for _, c := range complaints {
		date := c.Time().Format("20060102")
		counts[date] += 1
		if counts[date] > max {
			max = counts[date]
		}
	}
	if max == 0 {
		return nil
	}

	scale := &Scale{}
	frequency := make([]int, max)
	for _, n := range counts {
		frequency[n-1] += 1
		scale.Update(frequency[n-1])
	}

	io.WriteString(w, "Distribution of Complaints per day:\n")
	io.WriteString(w, scale.String())
	for freq, n := range frequency {
		if n == 0 {
			continue
		}
		fmt.Fprintf(w, "%2d complaints/day [ %3d days] %s\n", freq+1, n, strings.Repeat("∎", n/scale.Scale))
	}
	fmt.Fprint(w, "\n")
	return nil
}
