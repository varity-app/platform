package main

import (
	"time"
)

// generateDates generates a pair of arrays representing time intervals between two timestamps
func generateDates(start, end time.Time, minuteInterval int) ([]time.Time, []time.Time) {

	befores := []time.Time{}
	afters := []time.Time{}

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
						after := time.Date(year, month, day, hour, minute, 0, 0, time.UTC)
						before := after.Add(time.Duration(minuteInterval) * time.Minute)

						afters = append(afters, after)
						befores = append(befores, before)

						// fmt.Println(after.Format(time.RFC3339), before.Format(time.RFC3339))
					}
				}

			}
		}
	}

	return befores, afters
}

// Parse a date string
func parseDate(s string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}
