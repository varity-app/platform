package mentions

import (
	"time"

	b2i "github.com/varity-app/platform/scraping/internal/data/bigquery2influx"
	"github.com/varity-app/platform/scraping/internal/transforms/tickerext"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	write "github.com/influxdata/influxdb-client-go/v2/api/write"
)

const (
	cohortAll           = "all"
	measurementMentions = "mentions"

	tagSymbol    = "symbol"
	tagSubreddit = "subreddit"
	tagCohort    = "cohort"
	tagSource    = "source"

	fieldMentionsCount    = "mentions_count"
	fieldTargetedCount    = "targeted_count"
	fieldInquisitiveCount = "inquisitive_count"
)

// Aggregate aggregates a list of mentions by symbol, subreddit, cohort, and source.
func Aggregate(memberships []b2i.CohortMembership, mentions []b2i.Mention, ts time.Time) []*write.Point {

	aggs := aggregate(memberships, mentions)
	points := convertToPoints(aggs, ts)

	return points
}

// aggregate memberships and mentions
func aggregate(memberships []b2i.CohortMembership, mentions []b2i.Mention) map[[4]string][3]int {

	// Create a mapping for an author's id to its cohort memberships
	membershipsMap := make(map[string][]b2i.CohortMembership)
	for _, membership := range memberships {
		authorMemberships := membershipsMap[membership.AuthorID]

		membershipsMap[membership.AuthorID] = append(authorMemberships, membership)
	}

	// Aggregate mentions by cohort and subreddit
	aggs := make(map[[4]string][3]int)
	for _, mention := range mentions {
		memberships := membershipsMap[mention.AuthorID]

		// Always add a minimum membership to a fake "all" cohort.
		memberships = append(memberships, b2i.CohortMembership{Subreddit: mention.Subreddit, Cohort: cohortAll})

		// Skip mention if ticker in blacklist
		if contains(tickerext.TickerBlacklist, mention.Symbol) {
			continue
		}

		// Aggregate mentions by subreddit and cohort
		for _, membership := range memberships {
			if membership.Subreddit == mention.Subreddit {
				key := [4]string{mention.Symbol, mention.Subreddit, membership.Cohort, mention.Source}

				// Increment mention count
				curAgg := aggs[key]
				curAgg[0]++

				// Optionally increment targeted count
				if mention.Targeted {
					curAgg[1]++
				}

				// Optionally increment inquisitive count
				if mention.Inquisitive {
					curAgg[2]++
				}

				// Update map
				aggs[key] = curAgg
			}
		}
	}

	return aggs
}

// Convert aggregate map to a list of InfluxDB points
func convertToPoints(aggs map[[4]string][3]int, ts time.Time) []*write.Point {

	points := []*write.Point{}

	for key, val := range aggs {
		tags := map[string]string{
			tagSymbol:    key[0],
			tagSubreddit: key[1],
			tagCohort:    key[2],
			tagSource:    key[3],
		}

		fields := map[string]interface{}{
			fieldMentionsCount:    val[0],
			fieldTargetedCount:    val[1],
			fieldInquisitiveCount: val[2],
		}

		point := influxdb2.NewPoint(measurementMentions, tags, fields, ts)
		points = append(points, point)
	}

	return points
}

// Helper method for checking if element exists in array
func contains(arr []string, s string) bool {
	for _, el := range arr {
		if el == s {
			return true
		}
	}

	return false
}
