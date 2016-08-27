package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"db"

	"github.com/olekukonko/tablewriter"
)

func report() {
	reportByHour()
}

func reportByHour() {
	var hours [24]int
	var total int64
	complaints, err := db.Default.All()
	if err != nil {
		log.Fatalf("%s", err)
	}
	for _, c := range complaints {
		total += 1
		hours[c.Time().Hour()] += 1
	}
	start, stop := 0, 0
	for i, v := range hours {
		// if v > max {
		// 	max = v
		// }
		if (start == 0 || start == i-1) && v == 0 {
			start = i
		}
		if v != 0 {
			stop = i
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeader([]string{"Hour", "Complaints", ""})
	for h, v := range hours[start+1 : stop] {
		t := time.Date(2010, 1, 1, h+start, 0, 0, 0, time.UTC)
		table.Append([]string{t.Format("3pm"), fmt.Sprintf("%d", v), strings.Repeat("∎", v)})
	}
	table.SetFooter([]string{"Total:", fmt.Sprintf("%d", total), ""})
	table.Render()
}
