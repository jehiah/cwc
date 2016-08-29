package reporter

import (
	"io"

	"cwc/db"
)

func Run(d db.DB, w io.Writer) error {
	err := ByHour(d, w)
	if err != nil {
		return err
	}
	err = ByRegulation(d, w)
	if err != nil {
		return err
	}
	return nil
}
