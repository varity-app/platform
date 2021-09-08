package scrapers

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"

	"github.com/varity-app/platform/scraping/internal/common"
)

// Memory is an object that uses GCP Firestore as a key-value storage, and lets
// a scraper easily
type Memory struct {
	fsClient       *firestore.Client
	collectionName string
}

// MemoryOpts is a helper struct containing options passed to NewMemory
type MemoryOpts struct {
	CollectionName string
}

// NewMemory Initialize a new Memory
func NewMemory(ctx context.Context, opts MemoryOpts) (*Memory, error) {

	// Init firestore client
	fsClient, err := firestore.NewClient(ctx, common.GCPProjectID)
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %v", err)
	}

	return &Memory{
		fsClient:       fsClient,
		collectionName: opts.CollectionName,
	}, nil
}

// Close disconnects the memory connection
func (memory *Memory) Close() error {
	return memory.fsClient.Close()
}

// CheckNewItems takes two corresponding arrays of IDs and protobuf messages, and return the indexes of IDs that
// have not yet been saved in firestore.
func (memory *Memory) CheckNewItems(ctx context.Context, ids []string) ([]int, error) {
	// Create list of firestore refs for each submission
	refs := [](*firestore.DocumentRef){}
	collection := memory.fsClient.Collection(memory.collectionName)

	for _, id := range ids {
		refs = append(refs, collection.Doc(id))
	}

	// Fetch submissions from firestore
	docsnaps, err := memory.fsClient.GetAll(ctx, refs)
	if err != nil {
		return nil, fmt.Errorf("firestore.GetAll: %v", err)
	}

	// Check if each document has an existing reference.  If not, add it to the returned items
	var unseenIdxs []int

	for idx, snap := range docsnaps {
		// If there is no entry for an id,
		if !snap.Exists() {
			unseenIdxs = append(unseenIdxs, idx)
		}
	}

	return unseenIdxs, nil
}

//SaveItems Save entries for each specified ID in firestore.
func (memory *Memory) SaveItems(ctx context.Context, ids []string) error {
	collection := memory.fsClient.Collection(memory.collectionName)

	// Create batch of writes
	batch := memory.fsClient.Batch()
	willCommit := false

	// Check each snapshot for each document
	for _, id := range ids {
		batch.Create(collection.Doc(id), map[string]interface{}{})
		willCommit = true
	}

	// There are changes to write
	if willCommit {
		_, err := batch.Commit(ctx)
		if err != nil {
			return fmt.Errorf("firestore.WriteBatch: %v", err)
		}
	}

	return nil
}
