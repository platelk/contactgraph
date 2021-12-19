package userstore

import "testing"

func TestInMemoryImproved(t *testing.T) {
	runTestSuite(t, NewInMemoryImproved())
}
