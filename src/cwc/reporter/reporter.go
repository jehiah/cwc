package reporter

import (
	"encoding/json"
	"io"
	"log"

	"cwc/db"
	"cwc/reg"
)

func Run(d db.DB, w io.Writer) error {
	type reporter func(d db.DB, w io.Writer) error
	for _, r := range []reporter{ByHour, ByMonth, PerDay, ByRegulation, ByStatus, ByVehicle} {
		err := r(d, w)
		if err != nil {
			return err
		}
	}
	return nil
}

func JSON(d db.DB) ([]byte, error) {
	type record struct {
		Timestamp   int64     `json:"timestamp"`
		License     string    `json:"license_plate"`
		VehicleType string    `json:"vehicle_type"`
		Violations  []reg.Reg `json:"violations"`
		// Tweets []string `json:"tweets,omitempty"`
	}

	complaints, err := d.All()
	if err != nil {
		return nil, err
	}
	var data []*record
	for _, c := range complaints {
		record := &record{
			Timestamp:   c.Time().Unix(),
			License:     c.License(),
			VehicleType: "TAXI",
		}
		if ok, _ := d.ComplaintContains(c, " FHV "); ok {
			record.VehicleType = "FHV"
		}
		for _, r := range reg.All {
			if ok, _ := d.ComplaintContains(c, r.Code); ok {
				record.Violations = append(record.Violations, r)
			}
		}
		if len(record.Violations) == 0 {
			log.Printf("Warning: %s has no violations", c)
			continue
		}
		data = append(data, record)
	}
	body, err := json.Marshal(data)
	return body, err
}

func percent(n, total int) float32 {
	return (float32(n) / float32(total)) * 100
}
