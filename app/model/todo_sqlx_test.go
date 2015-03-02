package model

import (
	"log"
	"strconv"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func newTodoOrFatalSqlx(t *testing.T, title string) *Todo {
	var todo Todo
	todo.Title = title
	dbMap := getNewDbConnectionSqlx()
	if err := todo.SaveSqlx(dbMap); err != nil {
		t.Fatalf("new task: %v", err)
	}
	dbMap.DB.Close()
	return &todo
}

func TestNewTodoSqlx(t *testing.T) {
	title := "learn Go"
	todo := newTodoOrFatalSqlx(t, title)
	if todo.Title != title {
		t.Errorf("expected title %q, got %q", title, todo.Title)
	}
	if todo.IsCompleted {
		t.Errorf("new todo is done")
	}
}

func BenchmarkTodoInsertSqlx(b *testing.B) {
	for n := 0; n < b.N; n++ {
		dbMap := getNewDbConnectionSqlx()
		todo := Todo{Title: "abcdefghijklm"}
		if err := todo.SaveSqlx(dbMap); err != nil {
			dbMap.DB.Close()
			b.Errorf("Todo not saved", err)
		}
		dbMap.DB.Close()
	}
}

func BenchmarkTodoUpdateSqlx(b *testing.B) {
	for n := 0; n < b.N; n++ {
		dbMap := getNewDbConnectionSqlx()
		todo := Todo{TodoID: 10, Title: "abcdefghijklm"}
		if err := todo.UpdateSqlx(dbMap); err != nil {
			b.Errorf("Todo not updated", err)
		}
		dbMap.DB.Close()
	}
}

func BenchmarkTodoSelectAllSqlx(b *testing.B) {
	for n := 0; n < b.N; n++ {
		dbMap := getNewDbConnectionSqlx()
		_, _, err := GetAllTodosSqlx(dbMap, 0, 10000)
		if err != nil {
			b.Errorf("Todo not queried", err)
		}
		dbMap.DB.Close()
	}
	b.StopTimer()
	dbMap := getNewDbConnectionSqlx()
	DeleteAllSqlx(dbMap)
	dbMap.DB.Close()
	b.StartTimer()
}

func getNewDbConnectionSqlx() *sqlx.DB {
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
	dbMap, err := sqlx.Connect("mysql", strings.Join(dbStringSlice, ""))
	if err != nil {
		log.Fatalln(err, "sql.Open failed")
	}

	// db, err := sqlx.Connect("mysql", "user="+dbUserName+" dbname="+dbName+" ")

	return dbMap
}
