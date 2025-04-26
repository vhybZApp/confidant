package fact

// EffectHook runs a function when state changes
type EffectHook struct {
	callback func()
}

// NewEffectHook subscribes to a state change
func NewEffectHook[T any](state *ReactiveState[T], callback func(T)) *EffectHook {
	state.Subscribe(callback)
	return &EffectHook{callback: func() { callback(state.Get()) }}
}
