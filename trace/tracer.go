package trace

type Tracer interface {
	Tracee(...interface{})
}
