package prometheus

import (
	"cloudProject/pkg/cast"
	"sync/atomic"
)

type AtomicInt int64
type Stats struct {
	Errors    AtomicInt // show all errors (logical or internal)
	Successes AtomicInt
	Loads     AtomicInt // success + error
}

func (i *AtomicInt) Add(n int64) {
	atomic.AddInt64((*int64)(i), n)
}
func (i *AtomicInt) Set(n int64) {
	atomic.StoreInt64((*int64)(i), n)
}
func (i *AtomicInt) Get() int64 {
	return atomic.LoadInt64((*int64)(i))
}
func (s *Stats) AddError() {
	s.Errors.Add(1)
	s.AddLoads()
}
func (s *Stats) AddSuccess() {
	s.Successes.Add(1)
	s.AddLoads()
}
func (s *Stats) AddLoads() {
	s.Loads.Add(1)
}
func (s *Stats) GetError() float64 {
	toFloat64, _ := cast.ToFloat64(s.Errors.Get())
	return toFloat64
}
func (s *Stats) GetSuccess() float64 {
	toFloat64, _ := cast.ToFloat64(s.Successes.Get())
	return toFloat64
}
func (s *Stats) GetLoads() float64 {
	toFloat64, _ := cast.ToFloat64(s.Loads.Get())
	return toFloat64
}
