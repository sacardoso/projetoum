package main

import "testing"

func TestByteCountSI(t *testing.T) {
	type testStruct struct {
		num      int64
		expected string
	}

	tests := []testStruct{
		testStruct{num: int64(100), expected: "100 B"},
		testStruct{num: int64(2490), expected: "2.5 kB"},
		testStruct{num: int64(7200000), expected: "7.2 MB"},
		testStruct{num: int64(12800000000), expected: "12.8 GB"},
	}

	for _, test := range tests {
		result := ByteCountSI(test.num)
		if result != test.expected {
			t.Errorf("%s != %s", result, test.expected)
		}
	}
}
