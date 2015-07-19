package program

import (
	"testing"
)

func TestVersionCompare(t *testing.T) {
	vers := []VersionNo{
		{2, 2, 1}, {2, 2, 2}, {2, 2, 3},
		{3, 3, 2}, {3, 2, 2}, {3, 1, 2},
		{2, 3, 2}, {2, 2, 2}, {2, 1, 2},
		{1, 3, 2}, {1, 2, 2}, {1, 1, 2},
	}

	results := [3][12]int{
		[12]int{0, 0, 0, -1, -1, -1, -1, -1, 1, 1, 1, 1},
		[12]int{0, 0, 0, -1, -1, -1, -1, 0, 1, 1, 1, 1},
		[12]int{0, 0, 0, -1, -1, -1, -1, 1, 1, 1, 1, 1},
	}
	for j := 3; j < 12; j++ {
		for i := 0; i < 3; i++ {
			rs := vers[i].Compare(vers[j])
			t.Logf("compare(%s, %s) = %d", vers[i], vers[j], rs)
			if rs != results[i][j] {
				t.Errorf("compare(%s, %s) expects to get %d, but got %d", results[i][j], rs)
			}
		}
	}
}

func TestParseVersionNumber(t *testing.T) {
	kvs := map[string]VersionNo{
		"":       VersionNo{0, 0, 0},
		"1.":     VersionNo{1, 0, 0},
		"0.1":    VersionNo{0, 1, 0},
		"1.2.13": VersionNo{1, 2, 13},
		"7.-3.4": VersionNo{7, 0, 4},
	}
	for k, v := range kvs {
		if ver := ParseVersionNumber(k); ver.Compare(v) != 0 {
			t.Errorf("version number %q should be parsed to %v, but got %v", k, v, ver)
		}
	}
}
