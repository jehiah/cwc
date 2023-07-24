package exif

import (
	"path/filepath"
	"strings"
)

func ParseImageOrVideo(f string) (Exif, error) {
	ext := filepath.Ext(f)
	switch strings.ToLower(ext) {
	case ".jpeg", ".jpg", ".png":
		return ParseFile(f)
	case ".mov", ".mp4":
		return Ffmpeg(f)
	case ".heic":
		// TODO?
		// sips -g creation $file
		// creation: 2020:03:17 07:54:55
		return ParseFile(f)
	}
	return Exif{}, nil
}
