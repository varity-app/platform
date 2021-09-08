package prices

import (
	"context"
	"fmt"

	"github.com/varity-app/platform/scraping/internal/common"

	"cloud.google.com/go/bigquery"
)

// EODPriceSinkOpts is holds configuration options for the NewEODPriceSink constructor.
type EODPriceSinkOpts struct {
	DatasetName string
	TableName   string
}

// EODPriceSink is a BigQuery sink for storing EOD pricing data.
type EODPriceSink struct {
	bqClient    *bigquery.Client
	datasetName string
	tableName   string
}

// NewEODPriceSink constructs a new EODPriceSink
func NewEODPriceSink(ctx context.Context, opts EODPriceSinkOpts) (*EODPriceSink, error) {

	// Initialize bigquery client
	bqClient, err := bigquery.NewClient(ctx, common.GCPProjectID)
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %v", err)
	}

	return &EODPriceSink{
		bqClient:    bqClient,
		datasetName: opts.DatasetName,
		tableName:   opts.TableName,
	}, nil
}

// Close a sink's bigquery connection.
func (s *EODPriceSink) Close() error {
	return s.bqClient.Close()
}

// Sink saves a list of EOD prices to BigQuery
func (s *EODPriceSink) Sink(ctx context.Context, prices []EODPrice) error {
	// Create bigquery inserter
	inserter := s.bqClient.Dataset(s.datasetName).Table(s.tableName).Inserter()

	// Insert prices into bigquery table
	if err := inserter.Put(ctx, prices); err != nil {
		return fmt.Errorf("bigquery.Insert: %v", err)
	}

	return nil
}
