package main

import (
    "encoding/json"
    "net/http"
    "sync"
    "github.com/gorilla/handlers"
)

type Todo struct {
    ID   string `json:"id"`
    Task string `json:"task"`
}

var (
    todos []Todo
    mu    sync.Mutex
)

func getTodos(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
    var todo Todo
    // JSONデコードのエラーハンドリングを追加
    if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    mu.Lock()
    todos = append(todos, todo)
    mu.Unlock()
    w.WriteHeader(http.StatusCreated)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    mu.Lock()
    defer mu.Unlock()
    for i, todo := range todos {
        if todo.ID == id {
            todos = append(todos[:i], todos[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }
    // タスクが見つからない場合のレスポンスを改善
    http.Error(w, "Todo not found", http.StatusNotFound)
}

func main() {
    http.HandleFunc("/todos", getTodos)
    http.HandleFunc("/todos/create", createTodo)
    http.HandleFunc("/todos/delete", deleteTodo)

    // CORSを有効にする
    corsObj := handlers.AllowedOrigins([]string{"*"}) // すべてのオリジンを許可
    corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"}) // 許可するHTTPメソッド
    corsHeaders := handlers.AllowedHeaders([]string{"Content-Type"}) // 許可するヘッダー

    // サーバーを起動
    http.ListenAndServe(":8080", handlers.CORS(corsObj, corsMethods, corsHeaders)(http.DefaultServeMux))
}
