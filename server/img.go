package server

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/anthonynsimon/bild/transform"
	"github.com/jehiah/cwc/exif"
	"github.com/jehiah/cwc/internal/complaint"
)

func duplicateReader(r io.Reader) (io.Reader, io.Reader) {
	pr, pw := io.Pipe()
	tr := io.TeeReader(r, pw)
	return tr, pr
}

func (s *Server) Image(w http.ResponseWriter, r *http.Request, c complaint.Complaint, file string) {
	f, err := s.DB.OpenAttachment(c, file)
	if err != nil {
		http.Error(w, "INTERNAL_ERROR", 500)
		log.Printf("%s", err)
		return
	}
	defer f.Close()

	r1, r2 := duplicateReader(f)

	var x *exif.Exif
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		x, err = exif.Parse(r2)
		if err != nil {
			x = &exif.Exif{}
		}
	}()

	img, _, err := image.Decode(r1)
	if err != nil {
		http.Error(w, "INTERNAL_ERROR", 500)
		log.Printf("img decode error %s", err)
		return
	}
	wg.Wait()

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
