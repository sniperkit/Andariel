package task

type QueueEngine interface {
	FetchTasks(n int) ([]Task, error)
	DelTask(id interface{}) error
	ChangeActive(id interface{}, status int16) error
}
