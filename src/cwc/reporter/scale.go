package reporter

import (
	"fmt"
)

type Scale struct {
	Scale int
}

func (s *Scale) Update(n int) {
	v := (n / 70) + 1
	if v > s.Scale {
		s.Scale = v
	}
}
func (s Scale) String() string {
	if s.Scale <= 1 {
		return ""
	}
	return fmt.Sprintf("# each ∎ represents a count of %d\n", s.Scale)
}
