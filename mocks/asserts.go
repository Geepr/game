package mocks

import "testing"

func AssertDefault[T comparable](t *testing.T, value T) {
	var defaultValue T
	if value != defaultValue {
		t.Fail()
	}
}

func AssertCountEqual[T any](t *testing.T, value []T, expectedCount int) {
	if len(value) != expectedCount {
		t.Fail()
	}
}

func AssertArrayContains[T any](t *testing.T, collection []T, finder func(value T) bool) {
	for _, value := range collection {
		if finder(value) {
			return
		}
	}
	t.Fail()
}
