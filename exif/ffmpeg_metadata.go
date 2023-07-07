package exif

import (
	"bufio"
	"bytes"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Ffmpeg parses "exif" compatible data from a movie file using ffmpeg
func Ffmpeg(f string) (Exif, error) {
	e, err := getFFMetaData(f)
	if err != nil {
		return e, err
	}
	if e.Created.IsZero() {
		e.Created, err = getMovieCreationTime(f)
	}
	return e, err
}

func getFFMetaData(filePath string) (Exif, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-ffmetadata", "-i", filePath, "pipe:1")
	output, err := cmd.Output()
	if err != nil {
		return Exif{}, err
	}
	var e Exif
	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "=") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch key {
		case "location", "location-eng", "com.apple.quicktime.location.ISO6709":
			// https://developer.apple.com/library/archive/documentation/QuickTime/QTFF/Metadata/Metadata.html#//apple_ref/doc/uid/TP40000939-CH1-SW36:~:text=Group%20Video%20Music-,com.apple.quicktime.location.ISO6709,-%27mdta%E2%80%99
			// location=+40.7635-073.9853/
			// location-eng=+40.7635-073.9853/
			// Defined in ISO 6709:2008.
			// 	"+27.5916+086.5640+8850/"
			var remain string
			e.Lat, remain = parseISO6709Part(value)
			e.Long, _ = parseISO6709Part(remain)
		case "com.apple.quicktime.creationdate":
			e.Created, err = time.Parse("2006-01-02T15:04:05.000000Z", value)
			if err != nil {
				return e, err
			}
		}
	}

	return e, nil
}
func parseISO6709Part(s string) (f float64, remain string) {
	i := strings.Index(s, ".")
	split := strings.Index(s[i:], "+")
	if split == -1 {
		split = strings.Index(s[i:], "-")
	}
	f, _ = strconv.ParseFloat(s[:split+i], 64)
	return f, s[split+i:]
}

func getMovieCreationTime(filePath string) (time.Time, error) {
	// Use the appropriate method to extract the date-time metadata
	// from the .mov file. This can vary depending on the operating system and available tools.
	// Here's an example command using ffprobe (FFmpeg) on Unix-like systems:
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream_tags=creation_time", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
	output, err := cmd.Output()
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.Parse("2006-01-02T15:04:05.000000Z", strings.TrimSpace(string(output)))
	if err != nil {
		return t, err
	}
	nyc, _ := time.LoadLocation("America/New_York")
	return t.In(nyc), nil

	// use com.apple.quicktime.creationdate: 2023-01-13T16:34:26-05:00
	// which has timezone if available

}
