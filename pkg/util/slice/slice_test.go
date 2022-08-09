package slice

import (
	"testing"
)

func FuzzSlice(f *testing.F) {
	testcases := []string{"12345", "aaaa", "https://google.com", "sdhfskdlfh", "(*&%$(&$", ""}
	for _, tc := range testcases {
		f.Add(tc) // adding seed corpus
	}

	newSlice := []string{"aaaaa", "1234!", "234234234", "sdkfhsdfsd", "23423432"}

	f.Fuzz(func(t *testing.T, orig string) {
		newSlice = append(newSlice, orig)
		if !ContainsString(newSlice, orig) {
			t.Errorf("Expected to contain %s", orig)
		}
		newSlice = RemoveString(newSlice, orig)
		if ContainsString(newSlice, orig) {
			t.Errorf("Expected not to contain %s", orig)
		}
	})
}
