package input

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// UI is user-interface of input and output.
type UI struct {
	Reader io.Reader
	Writer io.Writer
}

var Default = &UI{
	Reader: os.Stdin,
	Writer: os.Stdout,
}

func (i *UI) readline() (string, error) {
	r := bufio.NewReader(i.Reader)
	s, err := r.ReadString('\n')
	if err == io.EOF {
		err = nil
	}
	return strings.TrimSpace(s), err
}
