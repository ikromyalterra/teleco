package threadpool

import (
	"time"

	"github.com/sepulsa/teleco/utils/logger"
)

// ThreadPool ...
type ThreadPool struct {
	Name           string
	MaxThread      int
	RoutineProcess chan struct{}
	Timeout        int
}

var (
	package_name_log     = "teleco/utils/threadpool"
	DefaultThreadNum     = 5
	DefaultThreadTimeout = 30
)

var thread map[string]*ThreadPool = make(map[string]*ThreadPool)

func (tp *ThreadPool) initPool() {
	tp.RoutineProcess = make(chan struct{}, tp.MaxThread)
	// Fill the dummy channel with maxThread empty struct.
	for i := 0; i < tp.MaxThread; i++ {
		tp.RoutineProcess <- struct{}{}
	}
	logger.Info().
		Str("event", "initPool").
		Str("package", package_name_log).
		Msgf("Init Thread Pool: %v, %v, %v", tp.Name, tp.MaxThread, tp.Timeout)
}

func Run(name string, maxThread int, timeout int, task Runnable) {
	tp, found := thread[name]
	if maxThread == 0 {
		maxThread = DefaultThreadNum
	}
	if timeout == 0 {
		timeout = DefaultThreadTimeout
	}
	if !found || tp.MaxThread != maxThread || tp.Timeout != timeout {
		threadPool := &ThreadPool{Name: name, MaxThread: maxThread, Timeout: timeout}
		threadPool.initPool()
		thread[name] = threadPool
		tp = thread[name]
	}

	done := make(chan bool)
	waitJobs := make(chan bool)

	go func() {
		<-done
		// Say that another goroutine can now start.
		tp.RoutineProcess <- struct{}{}
		waitJobs <- true
	}()

	// Try to receive from the routineProcess channel. When we have something,
	// it means we can start a new goroutine because another one finished.
	// Otherwise, it will block the execution until an execution spot is available.
	var isTimeout bool = false
	currentTime := time.Now()
	select {
	case <-tp.RoutineProcess:
		duration := currentTime.Add(time.Duration(tp.Timeout) * time.Second).Sub(time.Now())
		go func(t Runnable) {
			t.Run()
			if isTimeout {
				t.RunAfterTimeout()
			}
			done <- true
		}(task)
		select {
		case <-waitJobs:
			// Not do anything, just waiting channel signal sync process finished.
		case <-time.After(duration):
			task.RunWhenTimeout()
			isTimeout = true
		}
	case <-time.After(time.Duration(tp.Timeout) * time.Second):
		task.RunWhenFull()
	}
}
