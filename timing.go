package main

import "time"

// trackTime is a helper function to trackTime the time taken by a function.
// It takes a function as an argument and returns the result, error, and duration.
func trackTime[T any](f func() (T, error)) (T, error, time.Duration) {
	now := time.Now()
	result, err := f()
	duration := time.Since(now)
	return result, err, duration
}
