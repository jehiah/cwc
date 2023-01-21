package server

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"testing"
)

func TestDuplicateReader(t *testing.T) {
	v := "1234567890"
	b := strings.NewReader(v)
	r1, r2 := duplicateReader(b)

	var wg sync.WaitGroup

	for _, r := range []io.Reader{r1, r2} {
		r := r
		wg.Add(1)
		go func(r io.Reader) {
			defer wg.Done()
			body, err := io.ReadAll(r)
			if err != nil {
				t.Error(err.Error())
			}
			if !bytes.Equal(body, []byte(v)) {
				t.Errorf("got %q expected %q", string(body), v)
			}
		}(r)
		wg.Wait()
	}

}
