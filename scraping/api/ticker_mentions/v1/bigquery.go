package v1

import (
	"cloud.google.com/go/bigquery"
)

// Save implements the ValueSaver interface.  This is used for bigquery.
func (t *TickerMention) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"symbol":              t.Symbol,
		"parent_source":       t.ParentSource,
		"parent_id":           t.ParentId,
		"timestamp":           t.Timestamp.AsTime(),
		"symbol_counts":       t.SymbolCounts,
		"short_name_counts":   t.ShortNameCounts,
		"word_count":          t.WordCount,
		"question_mark_count": t.QuestionMarkCount,
	}, bigquery.NoDedupeID, nil
}
