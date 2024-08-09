package pipes

import (
	"testing"
)

func TestNewRing(t *testing.T) {
	// Test when elements are provided
	elements := []interface{}{1, 2, 3, 4, 5}
	ring, err := NewRing(elements)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ring == nil {
		t.Fatalf("expected ring to be non-nil")
	}
	if ring.size != len(elements) {
		t.Errorf("expected size %d, got %d", len(elements), ring.size)
	}

	if ring.idx != 0 {
		t.Errorf("expected index 0, got %d", ring.idx)
	}

	// Test when elements are empty
	ring, err = NewRing([]interface{}{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if ring != nil {
		t.Fatalf("expected ring to be nil")
	}
}

func TestRing_Next(t *testing.T) {
	elements := []interface{}{1, 2, 3, 4, 5}
	ring, err := NewRing(elements)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test the cycling through elements
	expectedSequence := []interface{}{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}
	for i, expected := range expectedSequence {
		if got := ring.Next(); got != expected {
			t.Errorf("at iteration %d, expected %v, got %v", i, expected, got)
		}
	}

	// Test the index wrapping
	for i := 0; i < len(elements)*2; i++ {
		ring.Next()
	}
	if ring.idx != 0 {
		t.Errorf("expected index to wrap around to 0, got %d", ring.idx)
	}

}
