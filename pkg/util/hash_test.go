package k8s

import (
	"testing"
)

func TestHash(t *testing.T) {
	if Hash("foo") != Hash("foo") {
		t.Errorf("Hash function result should be reproducible")
	}

	if Hash("foo") == Hash("bar") {
		t.Errorf("Hash function result should be unique if Namespace and Name do not match")
	}
}

func FuzzHash(f *testing.F) {
	testcases := []string{"foo", "bar", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}

	f.Fuzz(func(t *testing.T, orig string) {
		if Hash(orig) != Hash(orig) {
			t.Errorf("Hash function result should be reproducible")
		}
	})
}
