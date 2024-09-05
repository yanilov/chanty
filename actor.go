package chanty

type Actor[S any] struct {
	state  S
	system *ActorSystem
}

func NewActor[S any](state S, system *ActorSystem) *Actor[S] {
	actor := &Actor[S]{
		state:  state,
		system: system,
	}
	return actor
}

func Ask[S any, R any](actor *Actor[S], task func(*S) R) Future[R] {
	// adapt f(S)->R to f()->()
	fut := make(chan R, 1)
	adapter := func() {
		defer close(fut)
		state := &actor.state
		fut <- task(state)
	}

	// schedule the adapted task
	actor.system.schedule(adapter)

	return fut
}

func AskVoid[S any](actor *Actor[S], task func(*S)) Future[struct{}] {
	adapter := func(state *S) struct{} {
		task(state)
		return struct{}{}
	}
	return Ask(actor, adapter)
}

func Tell[S any](actor *Actor[S], task func(*S)) {
	adapter := func() {
		state := &actor.state
		task(state)
	}
	actor.system.schedule(adapter)
}

// Close closes the actor and returns the final state.
func Close[S any](actor *Actor[S]) *S {
	actor.system = nil
	return &actor.state
}
