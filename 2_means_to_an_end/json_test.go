package main

import (
	"testing"
)

func TestJsonDecoding(t *testing.T) {
	testcases := []struct {
		input string
		valid bool
	}{
		{`{}`, false},
		{`{"method": "isPrime", "number": 500}`, true},
		{`{"method": "isPrime"}`, false},
		{`{"number": 500}`, false},
		{`"number": 500}`, false},
	}

	for _, tc := range testcases {
		_, err := parseRequest([]byte(tc.input))
		if (err == nil) != tc.valid {
			t.Errorf("parseRequest(%v) should be %v", tc.input, tc.valid)
		}

	}

}
