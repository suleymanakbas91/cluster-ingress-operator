package k8s

import (
	"regexp"
	"testing"
)

// URI validates uri as being a valid http(s) uri and returns the uri scheme.
func TestURI(t *testing.T) {
	testCases := []struct {
		description string
		uri, scheme string
		expected    bool
	}{
		{
			description: "valid http uri with IP host and no port",
			uri:         "http://1.2.3.4",
			scheme:      "http",
			expected:    true,
		},
		{
			description: "valid http uri with IP host and backslash with no port",
			uri:         "http://1.2.3.4/",
			scheme:      "http",
			expected:    true,
		},
		{
			description: "valid http uri with IP host and port",
			uri:         "http://1.2.3.4:80",
			scheme:      "http",
			expected:    true,
		},
		{
			description: "valid http uri with IP host, port and backslash",
			uri:         "http://1.2.3.4:80/",
			scheme:      "http",
			expected:    true,
		},
		{
			description: "valid http uri with hostname",
			uri:         "http://redhat",
			scheme:      "http",
			expected:    true,
		},
		{
			description: "valid http uri with underscore in hostname",
			uri:         "http://red_hat.com",
			scheme:      "http",
			expected:    true,
		},
		{
			description: "valid http uri with FQDN",
			uri:         "http://www.redhat.com",
			scheme:      "http",
			expected:    true,
		},
		{
			description: "valid http uri with capitalized FQDN",
			uri:         "http://WWW.REDHAT.COM",
			scheme:      "http",
			expected:    true,
		},
		{
			description: "valid https uri with IP host and no port",
			uri:         "https://1.2.3.4",
			scheme:      "https",
			expected:    true,
		},
		{
			description: "valid https uri with mixed capitalization, port and backslash",
			uri:         "https://EXAMPLe.com:8080/",
			scheme:      "https",
			expected:    true,
		},
		{
			description: "http uri with invalid port number",
			uri:         "http://1.2.3.4:8080808080",
			scheme:      "http",
			expected:    false,
		},
		{
			description: "http uri with port number higher that the accepted range",
			uri:         "http://5.6.7.8:65536",
			scheme:      "http",
			expected:    false,
		},
		{
			description: "http uri with port number lower that the accepted range",
			uri:         "http://5.6.7.8:0",
			scheme:      "http",
			expected:    false,
		},
		{
			description: "missing uri scheme",
			uri:         "redhat.com",
			expected:    false,
		},
	}

	for _, tc := range testCases {
		scheme, err := URI(tc.uri)
		switch {
		case err != nil && tc.expected:
			t.Errorf("test %s failed: %v", tc.description, err)
		case err == nil && !tc.expected:
			t.Errorf("test %s expected to fail, but passed", tc.description)
		case err == nil && tc.expected:
			if scheme != tc.scheme {
				t.Errorf("unexpected scheme %s for test %s, expected scheme %s", scheme, tc.description, tc.scheme)
			}
		}
	}
}

func FuzzURI(f *testing.F) {
	testcases := []string{"http://1.2.3.4", "http://1.2.3.4/", "http://1.2.3.4:80",
		"http://1.2.3.4:80/", "http://redhat", "http://red_hat.com", "http://www.redhat.com",
		"http://WWW.REDHAT.COM", "https://1.2.3.4", "https://EXAMPLe.com:8080/"}
	for _, tc := range testcases {
		f.Add(tc) // adding seed corpus
	}

	f.Fuzz(func(t *testing.T, orig string) {
		pattern := `^(http[s]*:\/\/)([a-zA-Z\d\.]{2,})\.([a-zA-Z]{2,})(:1[0-9]{0,4}|:2[0-9]{0,4}|:3[0-9]{0,4}|:4[0-9]{0,4}|:5[0-9]{0,4}|:6[0-9]{0,3}[0-5])?$`
		re, err := regexp.Compile(pattern)
		if err != nil {
			t.Fatal("Failed to compile")
		}

		matched := re.Match([]byte(orig))
		if !matched {
			t.Skipf("Invalid URI: %s", orig)
		}

		// if orig is an int, skip
		// if orig is less than X characters, skip
		//
		scheme, err := URI(orig)
		if err != nil {
			t.Fatalf("Parsed scheme: %s", scheme)
		}
	})
}
