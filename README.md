[中文](README_CN.md)
# pooltracker
This is a component for tracking Sync.Pool usage in tests, allowing analysis of Sync.Pool usage without altering the original code.

## Features
- Tracks the number of calls to Pool.Get and Pool.Put.
- Track objects obtained from Pool.Get but not returned with Pool.Put.
- Tracks the value type of the Pool.

## Quick Start
Import the pooltracker package.
``` shell
go get github.com/jhue58/pooltracker
```
Use the Cover in your test cases.
``` go
type myStruct struct {
	val int
}

func TestTracker(t *testing.T) {
	pool := sync.Pool{New: func() any {
		return new(myStruct)
	}}

	Cover(t, func() {
		v:=pool.Get()
		pool.Put(v)
		pool.Get()
	})

}
```
Test Results
``` go
tracker.go:73: pool has leak value, Get 2, but Put 1,called by:
        tracker_test.go:20
```
You can also use NewTracker.
``` go
type myStruct struct {
	val int
}

func TestTracker(t *testing.T) {
	pool := sync.Pool{New: func() any {
		return new(myStruct)
	}}
	tracker := NewTracker()
	tracker.Track()
	defer tracker.Finish(t)
	
	v := pool.Get()
	vV := pool.Get()
	pool.Put(v)
	pool.Put(vV)
	pool.Put(false)
}
```
Test Results
``` go
tracker.go:82: pool value type should be *pooltracker.myStruct,but put bool,called by:
        tracker_test.go:24
```

## Note
Starting tracking will disable the original functionality of sync.Pool, leading to performance degradation. Avoid using it outside of test cases.