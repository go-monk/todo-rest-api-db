**Source: [https://github.com/go-monk/todo-rest-api-db](https://github.com/go-monk/todo-rest-api-db)**

This post builds on the [previous one](https://github.com/go-monk/todo-rest-api), where we created a simple REST API to manage an in-memory todo list. Keeping data in memory is fine for quick demos or tests, but real-world applications need to persist data so it survives server restarts. In this post, we’ll walk through how we added SQLite-based persistence to our todo API.

## The Problem: In-Memory Storage

Our original implementation used an in-memory store like this:

```go
type TaskStoreInMemory struct {
    sync.Mutex
    tasks  map[int]Task
    nextId int
}
```

This approach was fast and simple, but every time the server restarted, all tasks were lost. We needed a way to store tasks on disk.

Since we now anticipate having multiple implementations of the task store, let’s define an interface:

```go
type TaskStore interface {
    CreateTask(text string) (int, error)
    GetTasks() []Task
    GetTask(id int) (Task, error)
    DeleteTask(id int) error
}
```

The in-memory task store already almost satisfies this interface. We only need to modify its `CreateTask` method that now must return `(int, error)` instead of just `int`, allowing the handler to respond appropriately to database errors. We keep the in-memory implementation in `taskstore_mem.go`.

Now, let's look at something more persistent.

## The Solution: SQLite

SQLite is a lightweight, file-based database that’s perfect for small projects and prototyping. We implement a new `TaskStore` using SQLite, leveraging Go’s `database/sql` package and the popular `github.com/mattn/go-sqlite3` driver.

We create a new file, `taskstore_sqlite.go` and implement the TaskStore interface. Each method (CreateTask, GetTasks, GetTask, DeleteTask) is implemented using SQL queries to interact with the database.

## Switching Between In-Memory and SQLite Stores

We update our server to allow switching between the in-memory and SQLite-backed stores using a command-line flag. Here’s the new `main.go` logic:

```go
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
```

Now you can run the server with the `-persist` flag to use SQLite, or without it to use the in-memory store.


## Summary

Adding persistence to your API is a crucial step toward production readiness. With just a few changes, we upgraded our todo API from a toy example to a more robust service. SQLite is a great starting point, and the interface-based approach means you can later swap in more powerful databases with minimal changes. The new command-line flag makes it easy to switch between storage backends, and improved error handling ensures your API responds correctly to failures.