package quote

import (
	"testing"
)

func TestGenerateRandomId(t *testing.T) {
	length := 10
	id := generateRandomId(length)

	if len(id) != length {
		t.Errorf("Generated ID length is incorrect. Expected: %d, Got: %d", length, len(id))
	}
}
