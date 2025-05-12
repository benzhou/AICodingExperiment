package utils

import "time"

// TimeToMillis converts a Go time.Time to JavaScript-compatible milliseconds timestamp
func TimeToMillis(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UTC().UnixNano() / int64(time.Millisecond)
}

// EnsureTimestamps ensures that epoch time fields are populated from time.Time fields
// This can be used for any model that has CreatedAt/UpdatedAt and CreatedAtEpoch/UpdatedAtEpoch fields
func EnsureTimestamps(createdAt, updatedAt time.Time, createdAtEpoch, updatedAtEpoch *int64) {
	if *createdAtEpoch == 0 && !createdAt.IsZero() {
		*createdAtEpoch = TimeToMillis(createdAt)
	}

	if *updatedAtEpoch == 0 && !updatedAt.IsZero() {
		*updatedAtEpoch = TimeToMillis(updatedAt)
	}
}
