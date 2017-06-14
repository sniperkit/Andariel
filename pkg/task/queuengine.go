package task

type QueueEngine interface {
	FetchTasks(n int) ([]Task, error)
	DeleteTask(id interface{}) error
	Activate(id interface{}, status int16) error
}
