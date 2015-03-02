package model

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
)

func newTodoOrFatal(t *testing.T, title string) *Todo {
	var todo Todo
	todo.Title = title
	dbMap := getNewDbConnection()
	if err := todo.Save(dbMap); err != nil {
		t.Fatalf("new task: %v", err)
	}
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
		dbMap := getNewDbConnection()
		todo := Todo{Title: "abcdefghijklm"}
		if err := todo.Save(dbMap); err != nil {
			b.Errorf("Todo not updated", err)
		}
		dbMap.Db.Close()
	}
}

func BenchmarkTodoUpdateGorp(b *testing.B) {
	for n := 0; n < b.N; n++ {
		dbMap := getNewDbConnection()
		todo := Todo{TodoID: 10, Title: "abcdefghijklm"}
		if err := todo.Update(dbMap); err != nil {
			b.Errorf("Todo not queried", err)
		}
		dbMap.Db.Close()
	}
}

func BenchmarkTodoSelextAllGorp(b *testing.B) {
	for n := 0; n < b.N; n++ {
		dbMap := getNewDbConnection()
		_, _, err := GetAllTodos(dbMap, 0, 10000)
		if err != nil {
			b.Errorf("Todo not queried", err)
		}
		dbMap.Db.Close()
	}
	b.StopTimer()
	dbMap := getNewDbConnection()
	DeleteAll(dbMap)
	dbMap.Db.Close()
	b.StartTimer()
}

func getNewDbConnection() *gorp.DbMap {
	dbUserName := "root"
	dbPass := ""
	dbIp := "127.0.0.1"
	dbPortNo := 3306
	dbName := "wormhole"
	if dbPass != "" {
		dbPass = ":" + dbPass
	} else {
		dbPass = ":"
	}
	dbStringSlice := []string{dbUserName, dbPass, "@tcp(", dbIp, ":", strconv.Itoa(dbPortNo), ")/", dbName, "?parseTime=true"}
	db, err := sql.Open("mysql", strings.Join(dbStringSlice, ""))
	if err != nil {
		log.Fatalln(err, "sql.Open failed")
	}

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8"}}
	dbmap.AddTableWithName(Todo{}, "todos").SetKeys(true, "TodoID")

	return dbmap
}
