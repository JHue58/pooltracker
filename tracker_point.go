package pooltracker

import (
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
)

// Point is a pool tracker point
type Point struct {
	// Number of Get
	Get int32
	// Number of Put
	Put int32
	// The first Put value type of the pool
	Types reflect.Type
	// Different value type from Types
	// Map key is reflect.Type, value is []*runtime.Frames
	ITypes sync.Map
	// Leak value
	// Map key is value addr, value is *runtime.Frames
	LeakV sync.Map

	mu sync.Mutex
}

// getCalled is called when sync.Pool.Get() is called
func (t *Point) getCalled(v any, f *runtime.Frames) {

	t.LeakV.Store(v, f)
	atomic.AddInt32(&t.Get, 1)
}

// putCalled is called when sync.Pool.Put() is called
func (t *Point) putCalled(v any, f *runtime.Frames) {
	t.LeakV.Delete(v)
	tp := reflect.TypeOf(v)
	if atomic.AddInt32(&t.Put, 1) == 1 {
		t.Types = tp
		return
	}
	if t.Types != tp {
		t.mu.Lock()
		defer t.mu.Unlock()

		frames, _ := t.ITypes.LoadOrStore(tp, make([]*runtime.Frames, 0, 1))
		t.ITypes.Store(tp, append(frames.([]*runtime.Frames), f))
	}
}
