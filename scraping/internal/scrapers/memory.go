package scrapers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Memory is an object that uses Redis as a way to prevent duplicate scraped posts.
// It uses Redis as a database.
type Memory struct {
	rdb *redis.Client
}

// NewMemory Initialize a new Memory
func NewMemory(ctx context.Context, rdb *redis.Client) (*Memory, error) {
	return &Memory{
		rdb: rdb,
	}, nil
}

// CheckNewItems takes two corresponding arrays of IDs and protobuf messages, and return the indexes of IDs that
// have not yet been saved in firestore.
func (memory *Memory) CheckNewItems(ctx context.Context, prefix string, ids []string) ([]int, error) {

	// Check if each document has an existing reference.  If not, add it to the returned items
	var unseenIdxs []int

	pipe := memory.rdb.TxPipeline()
	defer pipe.Close()

	for _, id := range ids {
		key := fmt.Sprintf("%s:%s", prefix, id)
		pipe.Get(ctx, key)
	}

	// Execute pipe
	results, _ := pipe.Exec(ctx)
	for idx, result := range results {
		err := result.Err()
		if err == redis.Nil {
			unseenIdxs = append(unseenIdxs, idx)
		} else if err != nil {
			return nil, fmt.Errorf("redis.Get: %v", err)
		}
	}

	return unseenIdxs, nil
}

// SaveItems Save entries for each specified ID in firestore.
func (memory *Memory) SaveItems(ctx context.Context, prefix string, ids []string) error {
	pipe := memory.rdb.TxPipeline()
	defer pipe.Close()

	// Set entry in Redis for each ID
	for _, id := range ids {
		key := fmt.Sprintf("%s:%s", prefix, id)
		pipe.Set(ctx, key, "", 31*24*time.Hour).Err()
	}

	// Execute pipe
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("redis.Set: %v", err)
	}

	return nil
}
