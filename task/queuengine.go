package task

type QueueEngine interface {
	FetchTasks(n uint32) ([]Task, error)
	DelTask(id interface{}) error
}
