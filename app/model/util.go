package model

import (
	"log"

	"github.com/coopernurse/gorp"
	"github.com/mantishK/wormhole/storage"
)

func MysqlConnection() (dbMap *gorp.DbMap) {
	dbMap = storage.MysqlConnection()
	dbMap.AddTableWithName(Todo{}, "todos").SetKeys(true, "TodoID")

	//remove the following in production as it adds uneccessary complexity.
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		log.Fatalln("Error getting creating tables", err)
	}
	return
}
