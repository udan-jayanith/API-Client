package lazy_atomic

import (
	"sync"
)

type Value[T any] struct {
	// Use a mutex instead of a atomic.Value
	val   T
	Mutex sync.Mutex
}

func (v *Value[T]) LoadUnsafe() *T {
	return &v.val
}

func (v *Value[T]) Load() (val T) {
	v.Mutex.Lock()
	defer v.Mutex.Unlock()
	return v.val
}

func (v *Value[T]) Store(val T) {
	v.Mutex.Lock()
	defer v.Mutex.Unlock()
	v.val = val
}

func (v *Value[T]) Swap(new T) (old T) {
	v.Mutex.Lock()
	defer v.Mutex.Unlock()
	old = v.val
	v.val = new
	return old
}

func (v *Value[T]) Transaction(fn func(value T) (T, bool)) {
	v.Mutex.Lock()
	defer v.Mutex.Unlock()

	val, ok := fn(v.val)
	if ok {
		v.val = val
	}
}
