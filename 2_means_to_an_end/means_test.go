package main

import (
	"testing"
)

func TestIsPrime(t *testing.T) {
	testcases := []struct {
		num   float64
		prime bool
	}{
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{5.5, false},
		{27.0, false},
		{29, true},
		{30, false},
	}

	for _, tc := range testcases {
		prime := IsPrime(tc.num)
		if prime != tc.prime {
			t.Errorf("IsPrime(%v): %v, should be %v", tc.num, prime, tc.prime)
		}
	}
}
