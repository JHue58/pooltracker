package pooltracker

import (
	"github.com/agiledragon/gomonkey"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

type Tracker struct {
	p      *gomonkey.Patches
	enable int32
	m      sync.Map
}

func NewTracker() *Tracker {
	return new(Tracker)
}

// Track Start tracking,use Finish or Result to stop tracking.
func (t *Tracker) Track() {
	if !atomic.CompareAndSwapInt32(&t.enable, 0, 1) {
		return
	}
	t.p = gomonkey.NewPatches()
	sp := reflect.TypeOf(new(sync.Pool))
	t.p.ApplyMethod(sp, "Get", func(pool *sync.Pool) any {
		if atomic.LoadInt32(&t.enable) == 0 {
			return pool.New()
		}
		if pool == nil || pool.New == nil {
			return nil
		}
		newV := pool.New()
		val, _ := t.m.LoadOrStore(pool, new(Point))
		td := val.(*Point)
		pc := make([]uintptr, 512)
		n := runtime.Callers(2, pc)
		frames := runtime.CallersFrames(pc[:n])
		td.getCalled(newV, frames)
		return newV
	})

	t.p.ApplyMethod(sp, "Put", func(pool *sync.Pool, v any) {
		if atomic.LoadInt32(&t.enable) == 0 {
			return
		}
		if pool == nil || pool.New == nil {
			return
		}
		val, ok := t.m.Load(pool)
		if !ok {
			return
		}
		td := val.(*Point)
		pc := make([]uintptr, 512)
		n := runtime.Callers(2, pc)
		frames := runtime.CallersFrames(pc[:n])
		td.putCalled(v, frames)
	})

}

// Finish Stop tracking and get the test result.
func (t *Tracker) Finish(test *testing.T) {
	if atomic.LoadInt32(&t.enable) != 1 {
		return
	}
	defer t.reset()

	res := t.result()
	if res.HasLeak() > 0 {
		str := res.LeakStack()
		if str != "" {
			test.Fatalf(str)
		}

	}
	if res.HasInvalidType() > 0 {
		test.Fatalf(res.InvalidStack())
	}

}

// Result Stop tracking and get the tracked result.
func (t *Tracker) Result() (res TrackResult) {
	if atomic.LoadInt32(&t.enable) != 1 {
		return
	}
	defer t.reset()
	return t.result()
}

func (t *Tracker) result() (res TrackResult) {
	t.m.Range(func(key, value any) bool {
		res.Caught++
		td := value.(*Point)
		res.PoolTracked = append(res.PoolTracked, td)
		return true
	})
	return
}

func (t *Tracker) reset() {
	t.p.Reset()
	atomic.StoreInt32(&t.enable, 0)
	t.m = sync.Map{}
	t.p = nil
}
