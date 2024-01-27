package perf

type Client interface {
	Init() error
	Dial()
}
