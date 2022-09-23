package l27

import (
	"testing"
)

// Test that there are no deserialization errors or similar from loading all domains
func TestDomainGetAll(t *testing.T) {
	client := makeTestClient()

	client.Domains(CommonGetParams{Limit: 1000000})
}
