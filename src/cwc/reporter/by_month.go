package reporter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"cwc/db"

	"github.com/olekukonko/tablewriter"
)

func ByMonth(d db.DB, w io.Writer) error {
	counts := make(map[string]int)
	complaints, err := d.All()
	if err != nil {
		return err
	}
	for _, c := range complaints {
		month := c.Time().Format("200601")
		counts[month] += 1
	}
	var months []string
	for m, _ := range counts {
		months = append(months, m)
	}
	sort.Strings(months)

	// io.WriteString(w, "TLC Complaints by month\n")
	table := tablewriter.NewWriter(w)
	table.SetBorder(false)
	table.SetHeader([]string{"Month", "Complaints", ""})
	for _, month := range months {
		n := counts[month]
		table.Append([]string{month, fmt.Sprintf("%d", n), strings.Repeat("∎", n)})
	}
	table.Render()
	fmt.Fprint(w, "\n")
	return nil
}
