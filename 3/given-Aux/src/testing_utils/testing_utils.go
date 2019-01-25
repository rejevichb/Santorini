package testing_utils

import "testing"

func AssertTrue(result bool, t *testing.T) {
	if !result {
		t.Error("Expected true got false")
	}
}

func AssertEqual(t1 int, t2 int, t *testing.T) {
	if t1 != t2 {
		t.Errorf("Expected: %d, got: %d", t1, t2)
	}
}

func AssertEqualFloat(t1 float64, t2 float64, t *testing.T) {
	if t1 != t2 {
		t.Errorf("Expected: %f, got: %f", t1, t2)
	}
}
