package app

type Result[T any] struct {
	msg     string
	payload T
}

func (r Result[T]) Msg() string {
	return r.msg
}

func (r Result[T]) Payload() T {
	return r.payload
}

func ResultMsg[T any](msg string) Result[T] {
	return Result[T]{msg: msg}
}

func ResultPayload[T any](msg string, payload T) Result[T] {
	return Result[T]{msg: msg, payload: payload}
}

func ResultOK() Result[string] {
	return Result[string]{}
}
