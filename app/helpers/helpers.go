package helpers

import (
	"math"
	"os"
	"time"
)

// CreateFolder if it not exist
func CreateFolder(path string) error {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}

	return err
}

// GetDaysFromDate get days from creation date
func GetDaysFromDate(creationDate time.Time) int {
	hours := time.Since(creationDate).Hours()
	days := int(math.Ceil(hours / 24))

	return days
}
