package chanty

import (
	"context"
)

type Actor[S any] struct {
	state    S
	executor chan func(S)
}

func NewActor[S any](state S, buffer_len int) *Actor[S] {
	actor := &Actor[S]{
		state:    state,
		executor: make(chan func(S), buffer_len),
	}
	//event loop
	go func() {
		for task := range actor.executor {
			task(actor.state)
		}
	}()

	return actor
}

func Ask[S any, R any](actor *Actor[S], task func(state S) R) Response[R] {
	fut := make(chan R, 1)
	ctx, cancel := context.WithCancel(context.Background())

	adapter := func(s S) {
		fut <- task(s)
		close(fut)
	}
	go func() {
		select {
		case actor.executor <- adapter:
			break
		case <-ctx.Done():
			close(fut)
			break
		}
	}()
	return Response[R]{
		Future:      fut,
		Interaction: Interaction{cancel: cancel},
	}
}

func AskVoid[S any](actor *Actor[S], task func(state S)) Response[struct{}] {
	f := func(state S) struct{} {
		task(state)
		return struct{}{}
	}
	return Ask(actor, f)
}

func Tell[S any](actor *Actor[S], task func(state S)) Interaction {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case actor.executor <- task:
			break
		case <-ctx.Done():
			break
		}

	}()
	return Interaction{cancel: cancel}
}

// Close closes the actor and returns the final state.
func Close[S any](actor *Actor[S]) S {
	executor := actor.executor
	if executor != nil {
		actor.executor = nil
		close(executor)
		return actor.state
	}
	panic("actor already closed")
}
