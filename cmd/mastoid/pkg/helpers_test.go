package pkg

import "testing"

func TestExtractID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "numeric ID",
			input:    "117184748656417249",
			expected: "117184748656417249",
		},
		{
			name:     "normal Mastodon URL",
			input:    "https://mastodon.social/@user/117184748656417249",
			expected: "117184748656417249",
		},
		{
			name:     "remote account domain containing digit",
			input:    "https://mastodon.social/@stefano@rpi0w.stefanomarinelli.it/117184748656417249",
			expected: "117184748656417249",
		},
		{
			name:     "username containing digits",
			input:    "https://mastodon.social/@user123@example.social/117184748656417249",
			expected: "117184748656417249",
		},
		{
			name:     "query and fragment",
			input:    "https://mastodon.social/@user/117184748656417249?foo=123#456",
			expected: "117184748656417249",
		},
		{
			name:     "trailing slash",
			input:    "https://mastodon.social/@user/117184748656417249/",
			expected: "117184748656417249",
		},
		{
			name:     "non-numeric final path component",
			input:    "https://example.com/objects/ad3b097f-9599-4aa5-8ed5-ee463bbf7777",
			expected: "",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractID(tt.input); got != tt.expected {
				t.Fatalf("ExtractID(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
