package main

import (
	"time"

	b2i "github.com/varity-app/platform/scraping/internal/data/bigquery2influx"
	"github.com/varity-app/platform/scraping/internal/transforms"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	write "github.com/influxdata/influxdb-client-go/v2/api/write"
)

const (
	cohortAll           = "all"
	measurementMentions = "mentions"

	tagSymbol    = "symbol"
	tagSubreddit = "subreddit"
	tagCohort    = "cohort"

	fieldMentionsCount    = "mentions_count"
	fieldTargetedCount    = "targeted_count"
	fieldInquisitiveCount = "inquisitive_count"
)

// Aggregator aggregates a list of mentions by subreddit and cohort.
func aggregate(memberships []b2i.CohortMembership, mentions []b2i.Mention, ts time.Time) []*write.Point {

	// Create a mapping for an author's id to its cohort memberships
	membershipsMap := make(map[string][]b2i.CohortMembership)
	for _, membership := range memberships {
		authorMemberships := membershipsMap[membership.AuthorID]

		membershipsMap[membership.AuthorID] = append(authorMemberships, membership)
	}

	// Aggregate mentions by cohort and subreddit
	cohortCounts := make(map[[3]string][3]int)
	for _, mention := range mentions {
		memberships := membershipsMap[mention.AuthorID]

		// Always add a minimum membership to a fake "all" cohort.
		memberships = append(memberships, b2i.CohortMembership{Subreddit: mention.Subreddit, Cohort: cohortAll})

		// Skip mention if ticker in blacklist
		if contains(transforms.TickerBlacklist, mention.Symbol) {
			continue
		}

		// Aggregate mentions by subreddit and cohort
		for _, membership := range memberships {
			if membership.Subreddit == mention.Subreddit {
				key := [3]string{mention.Symbol, mention.Subreddit, membership.Cohort}

				// Increment mention count
				curAgg := cohortCounts[key]
				curAgg[0] += 1

				// Optionally increment targeted count
				if mention.Targeted {
					curAgg[1] += 1
				}

				// Optionally increment inquisitive count
				if mention.Inquisitve {
					curAgg[2] += 1
				}

				// Update map
				cohortCounts[key] = curAgg
			}
		}
	}

	points := convertToPoints(cohortCounts, ts)

	return points
}

// Convert aggregate map to a list of InfluxDB points
func convertToPoints(aggs map[[3]string][3]int, ts time.Time) []*write.Point {

	points := []*write.Point{}

	for key, val := range aggs {
		tags := map[string]string{
			tagSymbol:    key[0],
			tagSubreddit: key[1],
			tagCohort:    key[2],
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
	for _, el := range transforms.TickerBlacklist {
		if el == s {
			return true
		}
	}

	return false
}
