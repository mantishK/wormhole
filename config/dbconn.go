package config

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mantishK/wormhole/app/model"
)

func NewConnection() *gorp.DbMap {
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
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8"}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	// dbmap.AddTableWithName(model.Note{}, "note").SetKeys(true, "Note_id")
	dbmap.AddTableWithName(model.Todo{}, "todos").SetKeys(true, "TodoID")
	dbmap.AddTableWithName(model.User{}, "users").SetKeys(true, "UserID")
	dbmap.AddTableWithName(model.UserToken{}, "user_token")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}
func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
