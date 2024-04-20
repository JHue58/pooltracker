[English](README.md)
# pooltracker
这是一个用于在测试中跟踪Sync.Pool使用情况的组件，可以在不改变原代码的情况下分析Sync.Pool的使用情况。 

## 实现功能
- 跟踪Pool.Get和Pool.Put的调用次数
- 跟踪从Pool.Get获取的但未Put回去的对象
- 跟踪Pool的数据类型

## 快速开始
引入pooltracker包
``` shell
go get github.com/jhue58/pooltracker
```
在测试用例中使用Cover方法
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
测试结果
``` go
tracker.go:73: pool has leak value, Get 2, but Put 1,called by:
        tracker_test.go:20
```
也可以使用NewTracker
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
测试结果
``` go
tracker.go:82: pool value type should be *pooltracker.myStruct,but put bool,called by:
        tracker_test.go:24
```

## 注意
在开始跟踪后会使得sync.Pool的原本功能失效，导致性能降低，切勿在测试用例以外的地方使用