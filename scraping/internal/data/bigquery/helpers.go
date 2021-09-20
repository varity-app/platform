package bigquery

import (
	"fmt"

	"github.com/varity-app/platform/scraping/internal/common"
)

// Wrap a DBT BigQuery table name with the GCP project and dataset name.
func wrapDBTTable(deployment, tableName string) string {
	return fmt.Sprintf(
		"%s.%s_%s.%s",
		common.GCPProjectID,
		common.BigqueryDatasetDBTScraping,
		deployment,
		tableName,
	)
}

// Wrap a BigQuery scraping table name with the GCP project and dataset name.
func wrapTable(deployment, tableName string) string {
	return fmt.Sprintf(
		"%s.%s_%s.%s",
		common.GCPProjectID,
		common.BigqueryDatasetScraping,
		deployment,
		tableName,
	)
}
