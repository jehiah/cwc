package db

import (
	"strconv"
	"strings"
)

func findLongLatLine(lines []string) (long, lat float64) {
	// iterate backwards start with last log line
	for i := len(lines); i > 0; i-- {
		line := lines[i-1]
		if !strings.HasPrefix(line, "[ll:") {
			continue
		}
		line = line[4 : len(line)-1]
		chunks := strings.SplitN(line, ",", 2)
		if len(chunks) != 2 {
			continue
		}
		var err error
		long, err = strconv.ParseFloat(chunks[0], 64)
		if err != nil {
			continue
		}
		lat, err = strconv.ParseFloat(chunks[1], 64)
		if err != nil {
			continue
		}
		return
	}
	return
}
