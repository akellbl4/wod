package helpers

import (
	"os"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestCreateFolder(t *testing.T) {
	tmpDir := os.TempDir()
	err := CreateFolder(tmpDir)

	assert.NoError(t, err)

	err = CreateFolder(tmpDir)
	assert.NoError(t, err)
}

func TestGetDaysFromDate(t *testing.T) {
	var testSet = []struct {
		name   string
		date   time.Time
		result int
		patch  bool
	}{
		{
			name:   "today",
			date:   time.Now(),
			result: 1,
		},
		{
			name:   "unix zero",
			date:   time.Unix(0, 0),
			result: 18294,
			patch:  true,
		},
	}

	for _, tt := range testSet {
		t.Run(tt.name, func(t *testing.T) {
			if tt.patch {
				wayback := time.Date(2020, time.February, 2, 0, 0, 0, 0, time.UTC)
				patch := monkey.Patch(time.Now, func() time.Time { return wayback })
				defer patch.Unpatch()
			}
			assert.Equal(t, tt.result, GetDaysFromDate(tt.date))
		})
	}
}
