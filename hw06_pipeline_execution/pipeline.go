package hw06pipelineexecution

// In and Out are channel aliases used across the package.
type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func drain(c In) {
	for range c {
	}
}

func stageWorker(done In, stage Stage, in In) Out {
	if stage == nil {
		return in
	}

	stageOut := stage(in)
	bridge := make(Bi, 1)

	go func() {
		defer close(bridge)

		for {
			select {
			case <-done:
				go drain(stageOut)
				return

			case v, ok := <-stageOut:
				if !ok {
					return
				}

				select {
				case bridge <- v:
				case <-done:
					go drain(stageOut)
					return
				}
			}
		}
	}()

	return bridge
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	ch := in
	for _, st := range stages {
		ch = stageWorker(done, st, ch)
	}
	return ch
}
