package mocks

import "testing"

func AssertDefault[T comparable](t *testing.T, value T) {
	var defaultValue T
	if value != defaultValue {
		t.Fatalf("Expected default variable value")
	}
}

func AssertNotDefault[T comparable](t *testing.T, value T) {
	var defaultValue T
	if value == defaultValue {
		t.Fatalf("Expected not default variable value")
	}
}

func AssertCountEqual[T any](t *testing.T, value []T, expectedCount int) {
	if len(value) != expectedCount {
		t.Fatalf("Expected %d elements in collection, got %d instead", expectedCount, len(value))
	}
}

func AssertArrayContains[T any](t *testing.T, collection []T, finder func(value T) bool) {
	for _, value := range collection {
		if finder(value) {
			return
		}
	}
	t.Fatalf("Expected value was not found in collection")
}

func AssertEquals[T comparable](t *testing.T, actual T, expected T) {
	if actual != expected {
		t.Fatalf("Expected value doesn't match actual")
	}
}

func AssertEqualsNillable[T comparable](t *testing.T, actual *T, expected *T) {
	if (actual == nil && expected != nil) || (actual != nil && expected == nil) {
		t.Fatalf("Expected value doesn't match actual")
	}
	if actual == expected {
		return
	}
	if *actual != *expected {
		t.Fatalf("Expected value doesn't match actual")
	}
}

func CompareNillable[T comparable](left *T, right *T) bool {
	if left == nil {
		return right == nil
	}
	if right == nil {
		return left == nil
	}
	return *left == *right
}
