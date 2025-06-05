package main

import (
	"flag"
	"log"
	"net/http"
	"todo"

	"todo/handler"
)

func main() {
	persist := flag.Bool("persist", false, "persist tasks to SQLite")
	dbpath := flag.String("dbpath", "tasks.db", "path to SQLite file")
	flag.Parse()

	var store todo.TaskStore
	var err error
	if *persist {
		store, err = todo.NewSqliteTaskStore(*dbpath)
	} else {
		store, err = todo.NewInMemoryTaskStore()
	}
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	handler := handler.NewTaskHandler(store)
	mux.HandleFunc("POST /task", handler.AddTask)
	mux.HandleFunc("GET /tasks", handler.GetTasks)
	mux.HandleFunc("GET /task/{id}", handler.GetTask)
	mux.HandleFunc("DELETE /task/{id}", handler.DeleteTask)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
