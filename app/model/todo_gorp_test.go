package model

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func newTodoOrFatal(t *testing.T, title string) *Todo {
	var todo Todo
	todo.Title = title
	dbMap, _ := MysqlConnection()
	if err := todo.Save(dbMap); err != nil {
		t.Fatalf("new task: %v", err)
	}
	dbMap.Db.Close()
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
		dbMap, _ := MysqlConnection()
		if err := todo.Save(dbMap); err != nil {
			b.Errorf("Todo not inserted", err)
		}
		dbMap.Db.Close()
	}
}

func BenchmarkTodoUpdateGorp(b *testing.B) {
	for n := 0; n < b.N; n++ {
		todo := Todo{TodoID: 10, Title: "abcdefghijklm"}
		dbMap, _ := MysqlConnection()
		if err := todo.Update(dbMap); err != nil {
			b.Errorf("Todo not queried", err)
		}
		dbMap.Db.Close()
	}
}

func BenchmarkTodoSelectAllGorp(b *testing.B) {
	for n := 0; n < b.N; n++ {
		dbMap, _ := MysqlConnection()
		_, _, err := GetAllTodos(dbMap, 0, 10000)
		if err != nil {
			b.Errorf("Todo not queried", err)
		}
		dbMap.Db.Close()
	}
	b.StopTimer()
	dbMap, _ := MysqlConnection()
	DeleteAll(dbMap)
	dbMap.Db.Close()
	b.StartTimer()
}
