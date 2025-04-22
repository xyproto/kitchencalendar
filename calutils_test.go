package kitchencalendar

import (
	"testing"
	"time"
)

func TestWeekCalculation(t *testing.T) {
	// Test dates around New Year 2024/2025
	tests := []struct {
		date     time.Time
		expected int
	}{
		{time.Date(2024, 12, 28, 0, 0, 0, 0, time.UTC), 52}, // Week 52 of 2024
		{time.Date(2024, 12, 29, 0, 0, 0, 0, time.UTC), 52}, // Week 52 of 2024
		{time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC), 1},  // ISO Week 1 of 2025
		{time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC), 1},  // ISO Week 1 of 2025
		{time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), 1},    // ISO Week 1 of 2025
		{time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), 1},    // ISO Week 1 of 2025
		{time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC), 1},    // ISO Week 1 of 2025
	}

	for _, tt := range tests {
		got := GetWeekForDate(tt.date)
		if got != tt.expected {
			t.Errorf("GetWeekForDate(%v) = %d, want %d", tt.date.Format("2006-01-02"), got, tt.expected)
		}

		// Double-check with Go's built-in ISOWeek function
		_, week := tt.date.ISOWeek()
		if week != tt.expected {
			t.Errorf("ISOWeek(%v) = %d, which doesn't match our expected %d",
				tt.date.Format("2006-01-02"), week, tt.expected)
		}
	}
}

func TestFirstMondayOfWeek(t *testing.T) {
	tests := []struct {
		year     int
		week     int
		expected time.Time
	}{
		{2025, 1, time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC)},  // ISO Week 1 of 2025 starts on Dec 30, 2024
		{2025, 2, time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)},    // ISO Week 2 of 2025
		{2024, 52, time.Date(2024, 12, 23, 0, 0, 0, 0, time.UTC)}, // ISO Week 52 of 2024
	}

	for _, tt := range tests {
		got := FirstMondayOfWeek(tt.year, tt.week)
		if !got.Equal(tt.expected) {
			t.Errorf("FirstMondayOfWeek(%d, %d) = %v, want %v",
				tt.year, tt.week, got.Format("2006-01-02"), tt.expected.Format("2006-01-02"))
		}
	}
}

func TestFirstSundayOfWeek(t *testing.T) {
	tests := []struct {
		year     int
		week     int
		expected time.Time
	}{
		{2025, 1, time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)},    // Sunday of ISO Week 1 of 2025
		{2025, 2, time.Date(2025, 1, 12, 0, 0, 0, 0, time.UTC)},   // Sunday of ISO Week 2 of 2025
		{2024, 52, time.Date(2024, 12, 29, 0, 0, 0, 0, time.UTC)}, // Sunday of ISO Week 52 of 2024
	}

	for _, tt := range tests {
		got := FirstSundayOfWeek(tt.year, tt.week)
		if !got.Equal(tt.expected) {
			t.Errorf("FirstSundayOfWeek(%d, %d) = %v, want %v",
				tt.year, tt.week, got.Format("2006-01-02"), tt.expected.Format("2006-01-02"))
		}
	}
}
