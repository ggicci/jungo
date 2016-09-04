package request

import (
	"testing"
)

func TestAcceptEncodings(t *testing.T) {
	cps := map[string][]string{
		"":                                       []string{"identity"},
		"*":                                      []string{"*", "identity"},
		"compress,gzip":                          []string{"compress", "gzip", "identity"},
		"compress,gzip;q=1.0":                    []string{"compress", "gzip", "identity"},
		"gzip,deflate,sdch":                      []string{"gzip", "deflate", "sdch", "identity"},
		"gzip;q=0.3,deflate;q=0.9, identity;q=0": []string{"deflate", "gzip"},
		"gzip;q=1.0, identity;q=0.573, *;q=0":    []string{"gzip", "identity"},
		"gzip,*;q=0":                             []string{"gzip"},
		"*;q=0":                                  []string{},
	}

	for k, vs := range cps {
		result := acceptEncodings(k)
		if !equalStringSlice(vs, result) {
			t.Errorf("%q should get %v, got %v", k, vs, result)
		}
	}
}

func equalStringSlice(lhs, rhs []string) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for i := range lhs {
		if lhs[i] != rhs[i] {
			return false
		}
	}
	return true
}
