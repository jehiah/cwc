package complaint

import (
	"testing"
)

func TestFindTLCID(t *testing.T) {
	type testCase struct {
		line   string
		expect string
	}
	tests := []testCase{
		{"ion is needed. stip # 10073857s, mailed to d", "10073857"},
		{"ion is needed. stip10073857s, mailed to d", ""},
		{"14:56] complaint 10217074s.", "10217074"},
	}
	for _, tc := range tests {
		got := findTLCID([]string{tc.line})
		if got != tc.expect {
			t.Errorf("got %q expected %q for %q", got, tc.expect, tc.line)
		}
	}
}
