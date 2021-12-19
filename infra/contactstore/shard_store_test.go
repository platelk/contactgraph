package contactstore

import (
	"testing"
)

func TestNewShardStore(t *testing.T) {
	t.Run("1-shard", func(t *testing.T) {
		runTestSuiteContactStore(t, func() ContactStore {
			return NewShardStore(NewInMemoryMap())
		})
	})
	t.Run("3-shard", func(t *testing.T) {
		runTestSuiteContactStore(t, func() ContactStore {
			return NewShardStore(NewInMemoryMap(), NewInMemoryMap(), NewInMemoryMap())
		})
	})
}
