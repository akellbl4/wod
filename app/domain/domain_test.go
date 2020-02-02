package domain

import (
	"testing"
	"time"

	"bou.ke/monkey"
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
		name   string
		date   time.Time
		result int
		patch  bool
	}{
		{
			name: "today",
			date: time.Now(),
			result: 1,
		},
		{
			name: "unix zero",
			date: time.Unix(0, 0),
			result: 18294,
			patch: true,
		},
	}

	for _, tt := range testSet {
		t.Run(tt.name, func(t *testing.T) {
			if tt.patch {
				wayback := time.Date(2020, time.February, 2, 0, 0, 0, 0, time.UTC)
				patch := monkey.Patch(time.Now, func() time.Time { return wayback })
				defer patch.Unpatch()
			}
			assert.Equal(t, tt.result, GetDaysFromCreation(tt.date))
		})
	}
}
