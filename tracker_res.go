package pooltracker

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type TrackResult struct {
	// Number of tracked sync.Pool
	Caught int
	// All tracked sync.Pool
	PoolTracked []*Point
}

// HasLeak check sync.Pool has leak value.
func (t *TrackResult) HasLeak() (l int) {
	for _, point := range t.PoolTracked {
		if point.Put != point.Get {
			l++
			continue
		}
		point.LeakV.Range(func(_, _ any) bool {
			l++
			return false
		})
	}
	return
}

// HasInvalidType check sync.Pool has invalid type value.
func (t *TrackResult) HasInvalidType() (i int) {
	for _, point := range t.PoolTracked {
		point.ITypes.Range(func(_, _ any) bool {
			i++
			return false
		})
	}
	return
}

// LeakStack get all pool leaked value stack frame.
// It will only return the stack frame of the first leaked value for each pool.
func (t *TrackResult) LeakStack() string {
	builder := strings.Builder{}
	for _, point := range t.PoolTracked {
		point.LeakV.Range(func(_, value any) bool {
			builder.WriteString(fmt.Sprintf("pool has leak value, Get %d, but Put %d,called by:\n", point.Get, point.Put))
			frames := value.(*runtime.Frames)
			for {
				frame, ok := frames.Next()
				if !ok {
					break
				}
				builder.WriteString(fmt.Sprintf("function:%s,file:%s:%d\n", frame.Function, frame.File, frame.Line))
			}
			return false
		})

	}

	return builder.String()
}

// InvalidStack get all pool invalid type value stack frame.
func (t *TrackResult) InvalidStack() string {
	builder := strings.Builder{}
	for _, point := range t.PoolTracked {
		point.ITypes.Range(func(key, value any) bool {
			tp := key.(reflect.Type)
			fs := value.([]*runtime.Frames)
			builder.WriteString(fmt.Sprintf("pool value type should be %s,but put %s,called by:\n", point.Types, tp))
			frames := fs[0]
			for {
				frame, ok := frames.Next()
				if !ok {
					break
				}
				builder.WriteString(fmt.Sprintf("function:%s,file:%s:%d\n", frame.Function, frame.File, frame.Line))
			}
			return true
		})
	}
	return builder.String()
}
