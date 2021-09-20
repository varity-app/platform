package main

import (
	"time"
)

// generateMinuteRanges generates a pair of arrays (start, end) representing time intervals.
// Ranges are calculated down to the minute.
func generateMinuteRanges(start, end time.Time, minuteInterval int) ([]time.Time, []time.Time) {

	ends := []time.Time{}
	starts := []time.Time{}

	for year := start.Year(); year <= end.Year(); year++ {
		// Calculate max month
		maxMonth := time.Month(12)
		if year == end.Year() {
			maxMonth = end.Month()
		}

		// Iterate over each month
		for month := start.Month(); month <= maxMonth; month++ {
			// Calculate max day
			maxDay := 31
			if year == end.Year() && month == end.Month() {
				maxDay = end.Day()
			}

			for day := start.Day(); day <= maxDay; day++ {
				// Calculate max minute
				maxHour := 24
				if year == end.Year() && month == end.Month() && day == end.Day() {
					maxHour = end.Hour()
				}

				for hour := start.Hour(); hour <= maxHour; hour++ {
					// Calculate max minute
					maxMinute := 60
					if year == end.Year() && month == end.Month() && day == end.Day() && hour == end.Hour() {
						maxMinute = end.Minute()
					}

					// Iterate over each interval of minutes
					for minute := start.Minute(); minute <= maxMinute; minute += minuteInterval {
						start := time.Date(year, month, day, hour, minute, 0, 0, time.UTC)
						end := start.Add(time.Duration(minuteInterval) * time.Minute)

						starts = append(starts, start)
						ends = append(ends, end)
					}
				}

			}
		}
	}

	return starts, ends
}

// generateHourRanges generates a pair of arrays (start, end) representing time intervals.
// Ranges are calculated down to the hour.
func generateHourRanges(start, end time.Time) ([]time.Time, []time.Time) {

	ends := []time.Time{}
	starts := []time.Time{}

	for year := start.Year(); year <= end.Year(); year++ {
		// Calculate max month
		maxMonth := time.Month(12)
		if year == end.Year() {
			maxMonth = end.Month()
		}

		// Iterate over each month
		for month := start.Month(); month <= maxMonth; month++ {
			// Calculate max day
			maxDay := 31
			if year == end.Year() && month == end.Month() {
				maxDay = end.Day()
			}

			for day := start.Day(); day <= maxDay; day++ {
				// Calculate max minute
				maxHour := 24
				if year == end.Year() && month == end.Month() && day == end.Day() {
					maxHour = end.Hour()
				}

				for hour := start.Hour(); hour <= maxHour; hour++ {
					start := time.Date(year, month, day, hour, 0, 0, 0, time.UTC)
					end := start.Add(60 * time.Minute)

					starts = append(starts, start)
					ends = append(ends, end)
				}

			}
		}
	}

	return starts, ends
}

// Parse a date string
func parseDate(s string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}
