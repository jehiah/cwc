package reporter

import (
	"fmt"
	"io"
	"strings"

	"cwc/db"
)

func ByVehicle(d db.DB, w io.Writer) error {
	licenses := make(map[string][]db.Complaint)

	complaints, err := d.All()
	if err != nil {
		return err
	}

	var total, fhv, taxi int
	for _, c := range complaints {
		total++
		licenses[c.License()] = append(licenses[c.License()], c)
		if ok, _ := d.ComplaintContains(c, " FHV "); ok {
			fhv++
		} else {
			taxi++
		}
	}
	var preamble bool

	for l, cc := range licenses {
		if len(cc) < 2 {
			continue
		}
		if !preamble {
			fmt.Printf("License Plates w/ Multiple Reports:\n")
			preamble = true
		}
		var suffix []string
		for _, c := range cc {
			suffix = append(suffix, c.Time().Format("2006-01-02"))
		}
		fmt.Fprintf(w, "%-7s seen %d times (%s)\n", l, len(cc), strings.Join(suffix, ", "))
	}
	if preamble {
		fmt.Fprint(w, "\n")
	}

	totalLicenseCount := len(licenses)
	fmt.Fprintf(w, "Number of Unique License Plates: %d (of %d reports)\n", totalLicenseCount, total)
	fmt.Fprintf(w, "Taxi: %d (%0.1f%%) FHV: %d (%0.1f%%)\n", taxi, percent(taxi, totalLicenseCount), fhv, percent(fhv, totalLicenseCount))
	fmt.Fprint(w, "\n")

	return nil
}
