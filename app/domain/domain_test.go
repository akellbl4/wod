package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func timeMustBeParsed(date string) time.Time {
	t, err := time.Parse(time.RFC3339, date)

	if err != nil {
		panic(err)
	}

	return t
}

func TestGetDaysFromCreation(t *testing.T) {
	var testSet = []struct{
		name  string
		date  time.Time
		result int
	}{
		{
			name: "today",
			date: time.Now(),
			result: 1,
		},
		{
			name: "unix zero",
			date: timeMustBeParsed("1970-01-01T00:00:00.000Z"),
			result: 18294,
		},
	}

	for _, item := range testSet {
		t.Run(item.name, func(t *testing.T) {
			assert.Equal(t, item.result, GetDaysFromCreation(item.date))
		})
	}
}
