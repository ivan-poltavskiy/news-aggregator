package feed

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractDomainName(t *testing.T) {
	tests := []struct {
		url          string
		expectedName string
	}{
		{"https://www.example.com/path", "example"},
		{"http://example.com/path", "example"},
		{"https://example.com", "example"},
		{"https://sub.example.com/path", "sub"},
		{"invalid-url", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			domain := ExtractDomainName(tt.url)
			assert.Equal(t, tt.expectedName, domain, "Expected domain name to match")
		})
	}
}
