package main

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "sync"
    "log"
    "time" // 追加
    "github.com/gorilla/handlers"
    _ "github.com/lib/pq" // PostgreSQLドライバ
)

type Todo struct {
    ID      string `json:"id"`
    Task    string `json:"task"`
    DueDate string `json:"due_date"` // 期限を追加
}

var (
    db *sql.DB
    mu sync.Mutex
)

func init() {
    var err error
    // PostgreSQLに接続
    connStr := "host=db user=ss27914 password=meswSakeSaba39% dbname=todo_app sslmode=disable"

    for i := 0; i < 5; i++ { // 最大5回リトライ
        db, err = sql.Open("postgres", connStr)
        if err == nil {
            err = db.Ping() // 接続確認
        }
        if err == nil {
            break // 成功したらループを抜ける
        }
        log.Printf("Error connecting to the database: %v. Retrying in 2 seconds...", err)
        time.Sleep(2 * time.Second) // 2秒待機
    }

    if err != nil {
        log.Fatalf("Failed to connect to the database after retries: %v", err)
    }
    createTable()
}

func createTable() {
    createTableSQL := `CREATE TABLE IF NOT EXISTS todos (
        id TEXT PRIMARY KEY,
        task TEXT,
        due_date TEXT
    );`
    if _, err := db.Exec(createTableSQL); err != nil {
        panic(err)
    }
}

func getTodos(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()

    rows, err := db.Query("SELECT id, task, due_date FROM todos")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var todos []Todo
    for rows.Next() {
        var todo Todo
        if err := rows.Scan(&todo.ID, &todo.Task, &todo.DueDate); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        todos = append(todos, todo)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
    var todo Todo
    if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    mu.Lock()
    defer mu.Unlock()
    _, err := db.Exec("INSERT INTO todos (id, task, due_date) VALUES ($1, $2, $3)", todo.ID, todo.Task, todo.DueDate)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    mu.Lock()
    defer mu.Unlock()
    _, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func main() {
    http.HandleFunc("/todos", getTodos)
    http.HandleFunc("/todos/create", createTodo)
    http.HandleFunc("/todos/delete", deleteTodo)

    // CORSを有効にする
    corsObj := handlers.AllowedOrigins([]string{"*"})
    corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"})
    corsHeaders := handlers.AllowedHeaders([]string{"Content-Type"})

    // サーバーを起動
    log.Println("Starting server on :8080")
    http.ListenAndServe("0.0.0.0:8080", handlers.CORS(corsObj, corsMethods, corsHeaders)(http.DefaultServeMux))
}
