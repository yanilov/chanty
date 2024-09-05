package chanty

import (
	"context"
	"sync"
)

type ActorSystemConfig struct {
	NumWorkers   int
	Backpressure int
	Ctx          context.Context
}

func DefaultActorSystemConfig() *ActorSystemConfig {
	return &ActorSystemConfig{
		NumWorkers:   3,
		Backpressure: 1000,
		Ctx:          context.Background(),
	}
}

type ActorSystem struct {
	router    chan func()
	cancel    context.CancelFunc
	wgWorkers *sync.WaitGroup
}

func NewActorSystem(config ActorSystemConfig) *ActorSystem {

	ctx, cancel := context.WithCancel(config.Ctx)

	// initialize task router
	router := make(chan func(), config.Backpressure)

	var wg sync.WaitGroup

	// init workers' event loops
	for i := 0; i < config.NumWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case task, ok := <-router:
					if !ok {
						return
					}
					task()
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	return &ActorSystem{
		router,
		cancel,
		&wg,
	}
}

func (as *ActorSystem) Close() {
	// close the router s.t. new tasks are not accepted
	close(as.router)
	//signal all workers to stop
	as.cancel()
	// wait for graceful shutdown of all workers
	as.wgWorkers.Wait()
}

// schedule a task for execution. blocks the call if the actor system is busy
func (as *ActorSystem) schedule(task func()) {
	as.router <- task
}
