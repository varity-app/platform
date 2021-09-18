package bigquery2influx

import (
	"fmt"

	"github.com/varity-app/platform/scraping/internal/common"
)

// Wrap a DBT Bigquery table name with the GCP project and dataset name.
func wrapDBTTable(deployment, tableName string) string {
	return fmt.Sprintf(
		"%s.%s_%s.%s",
		common.GCPProjectID,
		common.BigqueryDatasetDBTScraping,
		deployment,
		tableName,
	)
}
