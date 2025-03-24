package mind

type Agent interface {
	Achieve(goal string, thread *Thread) error
}
