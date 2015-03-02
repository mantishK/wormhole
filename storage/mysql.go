package storage

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
)

var DbMap *gorp.DbMap

type settings struct {
	userName string
	pass     string
	ipAddr   string
	portNo   int
	dbName   string
}

func MysqlConnection() *gorp.DbMap {
	if DbMap != nil {
		return DbMap
	}

	s := mysqlSettings()
	if s.pass != "" {
		s.pass = ":" + s.pass
	} else {
		s.pass = ":"
	}
	sSlice := []string{s.userName, s.pass, "@tcp(", s.ipAddr, ":", strconv.Itoa(s.portNo), ")/", s.dbName, "?parseTime=true"}
	db, err := sql.Open("mysql", strings.Join(sSlice, ""))
	if err != nil {
		log.Fatalln("Error getting db connection", err)
	}

	dbMap := gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8"}}
	DbMap = &dbMap
	return DbMap
}

func CloseMysql() {
	if DbMap.Db != nil {
		DbMap.Db.Close()
		DbMap = nil
	}
}

func mysqlSettings() *settings {
	var s settings
	s.userName = "root"
	s.pass = ""
	s.ipAddr = "127.0.0.1"
	s.portNo = 3306
	s.dbName = "wormhole"
	return &s
}
