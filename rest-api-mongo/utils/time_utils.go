package utils

import (
	"time"
)

// GenerateID => generates a new todo ID
func GenerateID() int64 {
	now := time.Now()
	return now.Unix()
}
