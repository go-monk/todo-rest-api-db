package todo

type Task struct {
	Id   int
	Text string
}

type TaskStore interface {
	CreateTask(text string) (int, error)
	GetTasks() []Task
	GetTask(id int) (Task, error)
	DeleteTask(id int) error
}
