package fact

import (
	"sync"
)

// ReactiveState represents a reactive value with listeners
type ReactiveState[T any] struct {
	value     T
	listeners []func(T)
	mu        sync.Mutex
}

// NewReactiveState initializes a new reactive state
func NewReactiveState[T any](initialValue T) *ReactiveState[T] {
	return &ReactiveState[T]{value: initialValue}
}

// Get returns the current state value
func (s *ReactiveState[T]) Get() T {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.value
}

// Set updates the state and notifies all listeners
func (s *ReactiveState[T]) Set(newValue T) {
	s.mu.Lock()
	s.value = newValue
	listeners := s.listeners // Copy to avoid race conditions
	s.mu.Unlock()

	for _, listener := range listeners {
		listener(newValue)
	}
}

// Subscribe adds a listener that triggers when state changes
func (s *ReactiveState[T]) Subscribe(listener func(T)) {
	s.mu.Lock()
	s.listeners = append(s.listeners, listener)
	s.mu.Unlock()
}

// ReactiveListState manages a list of items reactively
type ReactiveListState[T any] struct {
	items     []T
	listeners []func([]T)
	mu        sync.Mutex
}

// NewReactiveListState initializes a reactive list
func NewReactiveListState[T any]() *ReactiveListState[T] {
	return &ReactiveListState[T]{items: []T{}}
}

// Get returns the list of items
func (rl *ReactiveListState[T]) Get() []T {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.items
}

// Set replaces the entire list and notifies subscribers
func (rl *ReactiveListState[T]) Set(newItems []T) {
	rl.mu.Lock()
	rl.items = newItems
	listeners := rl.listeners
	rl.mu.Unlock()

	for _, listener := range listeners {
		listener(newItems)
	}
}

// Add appends a new item and notifies subscribers
func (rl *ReactiveListState[T]) Add(item T) {
	rl.mu.Lock()
	rl.items = append(rl.items, item)
	listeners := rl.listeners
	rl.mu.Unlock()

	for _, listener := range listeners {
		listener(rl.items)
	}
}

// Remove deletes an item at an index and notifies subscribers
func (rl *ReactiveListState[T]) Remove(index int) {
	rl.mu.Lock()
	if index < 0 || index >= len(rl.items) {
		rl.mu.Unlock()
		return
	}
	rl.items = append(rl.items[:index], rl.items[index+1:]...)
	listeners := rl.listeners
	rl.mu.Unlock()

	for _, listener := range listeners {
		listener(rl.items)
	}
}

// Subscribe adds a listener to watch list changes
func (rl *ReactiveListState[T]) Subscribe(listener func([]T)) {
	rl.mu.Lock()
	rl.listeners = append(rl.listeners, listener)
	rl.mu.Unlock()
}

// AsyncState manages asynchronous data fetching
type AsyncState[T any] struct {
	value     T
	loading   bool
	error     error
	listeners []func(T, bool, error)
	mu        sync.Mutex
}

// NewAsyncState creates an async state
func NewAsyncState[T any](fetchFunc func() (T, error)) *AsyncState[T] {
	state := &AsyncState[T]{loading: true}
	go state.load(fetchFunc)
	return state
}

// Get returns the current value, loading state, and error
func (as *AsyncState[T]) Get() (T, bool, error) {
	as.mu.Lock()
	defer as.mu.Unlock()
	return as.value, as.loading, as.error
}

// Subscribe listens for updates
func (as *AsyncState[T]) Subscribe(listener func(T, bool, error)) {
	as.mu.Lock()
	as.listeners = append(as.listeners, listener)
	as.mu.Unlock()
}

// load fetches data asynchronously
func (as *AsyncState[T]) load(fetchFunc func() (T, error)) {
	result, err := fetchFunc()

	as.mu.Lock()
	as.value = result
	as.error = err
	as.loading = false
	listeners := as.listeners
	as.mu.Unlock()

	// Notify all listeners
	for _, listener := range listeners {
		listener(result, false, err)
	}
}
