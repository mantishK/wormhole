package model

import (
	"log"
	"strconv"
	"strings"
	"testing"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func newTodoOrFatalSql(t *testing.T, title string) *Todo {
	var todo Todo
	todo.Title = title
	dbMap := getNewDbConnectionSql()
	if err := todo.SaveSql(dbMap); err != nil {
		t.Fatalf("new task: %v", err)
	}
	dbMap.Close()
	return &todo
}

func TestNewTodoSql(t *testing.T) {
	title := "learn Go"
	todo := newTodoOrFatalSql(t, title)
	if todo.Title != title {
		t.Errorf("expected title %q, got %q", title, todo.Title)
	}
	if todo.IsCompleted {
		t.Errorf("new todo is done")
	}
}

func BenchmarkTodoInsertSql(b *testing.B) {
	for n := 0; n < b.N; n++ {
		dbMap := getNewDbConnectionSql()
		todo := Todo{Title: "abcdefghijklm"}
		if err := todo.SaveSql(dbMap); err != nil {
			b.Errorf("Todo not saved", err)
		}
		dbMap.Close()
	}
}

func BenchmarkTodoUpdateSql(b *testing.B) {
	for n := 0; n < b.N; n++ {
		dbMap := getNewDbConnectionSql()
		todo := Todo{TodoID: 10, Title: "abcdefghijklm"}
		if err := todo.UpdateSql(dbMap); err != nil {
			b.Errorf("Todo not updated", err)
		}
		dbMap.Close()
	}
}

func BenchmarkTodoSelectAllSql(b *testing.B) {
	for n := 0; n < b.N; n++ {
		dbMap := getNewDbConnectionSql()
		_, _, err := GetAllTodosSql(dbMap, 0, 10000)
		if err != nil {
			b.Errorf("Todo not queried", err)
		}
		dbMap.Close()
	}
	b.StopTimer()
	dbMap := getNewDbConnectionSql()
	DeleteAllSql(dbMap)
	dbMap.Close()
	b.StartTimer()
}

func getNewDbConnectionSql() *sql.DB {
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
	dbMap, err := sql.Open("mysql", strings.Join(dbStringSlice, ""))
	if err != nil {
		log.Fatalln(err, "sql.Open failed")
	}

	// db, err := sqlx.Connect("mysql", "user="+dbUserName+" dbname="+dbName+" ")

	return dbMap
}
