package contactstore

import (
	"testing"
)

func TestNewInMemoryMap(t *testing.T) {
	runTestSuiteContactStore(t, func() ContactStore {
		return NewInMemoryMap()
	})
}
