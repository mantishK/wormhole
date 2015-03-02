package model

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mantishK/wormhole/storage"
)

func newTodoOrFatal(t *testing.T, title string) *Todo {
	var todo Todo
	todo.Title = title
	if err := todo.Save(); err != nil {
		t.Fatalf("new task: %v", err)
	}
	storage.CloseMysql()
	return &todo
}

func TestNewTodo(t *testing.T) {
	title := "learn Go"
	todo := newTodoOrFatal(t, title)
	if todo.Title != title {
		t.Errorf("expected title %q, got %q", title, todo.Title)
	}
	if todo.IsCompleted {
		t.Errorf("new todo is done")
	}
}

func BenchmarkTodoInsertGorp(b *testing.B) {
	for n := 0; n < b.N; n++ {
		todo := Todo{Title: "abcdefghijklm"}
		if err := todo.Save(); err != nil {
			b.Errorf("Todo not inserted", err)
		}
		storage.CloseMysql()
	}
}

func BenchmarkTodoUpdateGorp(b *testing.B) {
	for n := 0; n < b.N; n++ {
		todo := Todo{TodoID: 10, Title: "abcdefghijklm"}
		if err := todo.Update(); err != nil {
			b.Errorf("Todo not queried", err)
		}
		storage.CloseMysql()
	}
}

func BenchmarkTodoSelextAllGorp(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _, err := GetAllTodos(0, 10000)
		if err != nil {
			b.Errorf("Todo not queried", err)
		}
		storage.CloseMysql()
	}
	b.StopTimer()
	DeleteAll()
	storage.CloseMysql()
	b.StartTimer()
}
