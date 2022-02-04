package consumer

// Runnable is interface for the jobs that will be executed by the worker consumer
type Runnable interface {
	Run(payload string)
}
