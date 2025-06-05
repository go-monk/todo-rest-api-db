package todo

import (
	"fmt"
	"sync"
)

type TaskStoreInMemory struct {
	sync.Mutex

	tasks  map[int]Task
	nextId int
}

func NewInMemoryTaskStore() (*TaskStoreInMemory, error) {
	ts := &TaskStoreInMemory{}
	ts.tasks = make(map[int]Task)
	ts.nextId = 0
	return ts, nil
}

func (ts *TaskStoreInMemory) CreateTask(text string) (int, error) {
	ts.Lock()
	defer ts.Unlock()

	task := Task{
		Id:   ts.nextId,
		Text: text,
	}

	ts.tasks[ts.nextId] = task
	ts.nextId++
	return task.Id, nil
}

func (ts *TaskStoreInMemory) GetTasks() []Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []Task
	for _, t := range ts.tasks {
		tasks = append(tasks, t)
	}
	return tasks
}

func (ts *TaskStoreInMemory) GetTask(id int) (Task, error) {
	ts.Lock()
	defer ts.Unlock()

	t, ok := ts.tasks[id]
	if !ok {
		return Task{}, fmt.Errorf("task with id=%d not found", id)
	}
	return t, nil
}

func (ts *TaskStoreInMemory) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()

	if _, ok := ts.tasks[id]; !ok {
		return fmt.Errorf("task with id=%d not found", id)
	}

	delete(ts.tasks, id)
	return nil
}
