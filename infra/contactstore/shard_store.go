package contactstore

import (
	"context"
	"sync"

	"github.com/platelk/contactgraph/domain/models/contacts"
	"github.com/platelk/contactgraph/domain/models/users"
)

type ShardStore struct {
	shards []IterableStore
}

func NewShardStore(shards ...IterableStore) *ShardStore {
	return &ShardStore{shards: shards}
}

func (s *ShardStore) Connect(ctx context.Context, from, to users.ID) error {
	return s.shards[int(from)%len(s.shards)].Connect(ctx, from, to)
}

func (s *ShardStore) Lookup(ctx context.Context, userID users.ID) (contacts.List, error) {
	return s.shards[int(userID)%len(s.shards)].Lookup(ctx, userID)
}

func (s *ShardStore) ReverseLookup(ctx context.Context, userID users.ID) (contacts.List, error) {
	var wg, rwg sync.WaitGroup
	result := contacts.List{}
	resultChan := make(chan users.ID)
	rwg.Add(1)
	go func() {
		for friend := range resultChan {
			result[friend] = struct{}{}
		}
		rwg.Done()
	}()
	for _, shard := range s.shards {
		wg.Add(1)
		go func(shard IterableStore) {
			_ = shard.Iterate(ctx, func(friend users.ID, friends contacts.List) {
				if _, ok := friends[userID]; ok {
					resultChan <- friend
				}
			})
			wg.Done()
		}(shard)
	}
	wg.Wait()
	close(resultChan)
	rwg.Wait()
	return result, nil
}

func (s *ShardStore) Len() uint {
	var total uint
	for _, shard := range s.shards {
		total += shard.Len()
	}
	return total
}
