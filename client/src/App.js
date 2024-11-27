import React, { useState, useEffect } from 'react';
import './App.css';

function App() {
  const [todos, setTodos] = useState([]);
  const [task, setTask] = useState('');
  const [error, setError] = useState(''); // エラーメッセージ用の状態を追加

  const fetchTodos = async () => {
    try {
      const response = await fetch('http://localhost:8080/todos');
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      const data = await response.json();
      setTodos(data || []);
    } catch (error) {
      console.error('Error fetching todos:', error);
      setTodos([]);
    }
  };

  const addTodo = async () => {
    if (!task) return;

    const newTodo = { id: Date.now().toString(), task };
    setTodos((prevTodos) => [...prevTodos, newTodo]); // 新しいタスクを即座に追加

    try {
      await fetch('http://localhost:8080/todos/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newTodo),
      });
      setTask('');
      setError(''); // エラーをリセット
      fetchTodos();
    } catch (error) {
      console.error('Error adding todo:', error);
      setError('タスクの追加に失敗しました。'); // エラーメッセージを設定
    }
  };

  const deleteTodo = async (id) => {
    try {
      await fetch(`http://localhost:8080/todos/delete?id=${id}`, {
        method: 'DELETE',
      });
      setTodos((prevTodos) => prevTodos.filter(todo => todo.id !== id)); // 削除したタスクを即座にリストから削除
    } catch (error) {
      console.error('Error deleting todo:', error);
      setError('タスクの削除に失敗しました。'); // エラーメッセージを設定
    }
  };

  useEffect(() => {
    fetchTodos();
  }, []);

  return (
    <div className="App">
      <h1>ToDo List</h1>
      {error && <p style={{ color: 'red' }}>{error}</p>} {/* エラーメッセージを表示 */}
      <input
        type="text"
        value={task}
        onChange={(e) => setTask(e.target.value)}
        placeholder="Add a new task"
      />
      <button onClick={addTodo}>Add</button>
      <ul>
        {Array.isArray(todos) && todos.map((todo) => (
          <li key={todo.id}>
            {todo.task}
            <button onClick={() => deleteTodo(todo.id)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
