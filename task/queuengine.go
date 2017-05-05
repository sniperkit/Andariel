package task

type QueueEngine interface {
	FetchTask() (*Task, error)
	DelTask(id interface{}) error
}
