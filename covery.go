package pooltracker

import (
	"testing"
)

var (
	defaultTracker = new(Tracker)
)

// Cover your test case.
func Cover(t *testing.T, f func()) {
	defaultTracker.Track()
	defer defaultTracker.Finish(t)
	f()
}
