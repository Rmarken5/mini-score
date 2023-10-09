package general

import "time"

func StartTime(t time.Time) time.Time {

	// Calculate the number of days to subtract to reach the previous Tuesday
	daysToSubtract := int(t.Weekday() - time.Tuesday)
	if daysToSubtract < 0 {
		daysToSubtract += 7 // Add 7 days to move to the previous week
	}

	// Subtract the days to reach the previous Tuesday
	previousTuesday := t.AddDate(0, 0, -daysToSubtract)

	// Set the time to midnight (00:00:00)
	start := time.Date(previousTuesday.Year(), previousTuesday.Month(), previousTuesday.Day(), 0, 0, 0, 0, time.Local)

	return start.UTC()
}

func EndTime(t time.Time) time.Time {

	// Calculate the number of days to add to reach the coming Monday
	daysToAdd := int(time.Monday - t.Weekday())
	if daysToAdd < 0 {
		daysToAdd += 7 // Add 7 days to move to the next week
	}

	// Add the days to reach the coming Monday
	comingMonday := t.AddDate(0, 0, daysToAdd)

	// Set the time to 11:59 PM
	end := time.Date(comingMonday.Year(), comingMonday.Month(), comingMonday.Day(), 23, 59, 59, 0, time.Local)

	return end.UTC()
}
