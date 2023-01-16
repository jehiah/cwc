package server

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/anthonynsimon/bild/transform"
	"github.com/jehiah/cwc/exif"
	"github.com/jehiah/cwc/internal/complaint"
)

func (s *Server) Image(w http.ResponseWriter, r *http.Request, c complaint.Complaint, file string) {
	path := s.DB.FullPath(c)
	file = filepath.Join(path, file)

	f, err := os.Open(file)
	if err != nil {
		http.Error(w, "INTERNAL_ERROR", 500)
		log.Printf("%s", err)
		return
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		http.Error(w, "INTERNAL_ERROR", 500)
		log.Printf("img decode error %s", err)
		return
	}

	x, err := exif.Parse(file)
	if err != nil {
		x = &exif.Exif{}
	}

	// rotate & transform
	r.ParseForm()
	width := img.Bounds().Dx()
	// estimate rotation
	switch x.ExifRotation {
	case 90, 270:
		width = img.Bounds().Dy()
	}
	origWidth := width

	if v := r.Form.Get("w"); v != "" {
		if w, err := strconv.Atoi(v); err == nil && w > 0 && w < width {
			width = w
		}
	}

	// double for @2x if possible
	if width*2 < origWidth {
		width = width * 2
	}

	// rotate & flip
	if x.ExifFlip {
		img = transform.FlipH(img)
	}
	if x.ExifRotation > 0 {
		img = transform.Rotate(img, x.ExifRotation, &transform.RotationOptions{ResizeBounds: true})
	}
	// resize
	if width != origWidth {
		ratio := float64(width) / float64(origWidth)
		img = transform.Resize(img, width, int(float64(img.Bounds().Dy())*ratio), transform.NearestNeighbor)
	}

	// now encode
	w.Header().Set("Content-Type", "image/png")
	err = png.Encode(w, img)
	if err != nil {
		http.Error(w, "INTERNAL_ERROR", 500)
		log.Printf("%s", err)
	}
}

func (s *Server) Download(w http.ResponseWriter, r *http.Request, c complaint.Complaint, file string) {
	path := s.DB.FullPath(c)
	staticServer := http.StripPrefix(fmt.Sprintf(s.BasePath+"complaint/%s/", c.ID()), http.FileServer(http.Dir(path)))
	staticServer.ServeHTTP(w, r)
}
