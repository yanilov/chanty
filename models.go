package chanty

import (
	"context"
)

type Interaction struct {
	cancel context.CancelFunc
}

func (i *Interaction) Cancel() {
	i.cancel()
}

type Future[T any] <-chan T

type Response[T any] struct {
	Future      Future[T]
	Interaction Interaction
}
