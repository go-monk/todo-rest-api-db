package handler

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"strconv"
	"todo"
)

type taskHandler struct {
	store todo.TaskStore
}

func NewTaskHandler(store todo.TaskStore) *taskHandler {
	return &taskHandler{store: store}
}

func (th *taskHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling add task at %s\n", r.URL.Path)

	// Types used internally in this handler to (de-)serialize
	// the request and response from/to JSON.
	type RequestTask struct {
		Text string `json:"text"`
	}
	type ResponseId struct {
		Id int `json:"id"`
	}

	// Enforce a JSON Content-Type.
	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	var t RequestTask
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := th.store.CreateTask(t.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, ResponseId{Id: id})
}

func (th *taskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get all tasks at %s\n", r.URL.Path)

	tasks := th.store.GetTasks()
	renderJSON(w, tasks)
}

func (th *taskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling delete task at %s\n", r.URL.Path)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = th.store.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (th *taskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get task at %s\n", r.URL.Path)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	t, err := th.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, t)
}
