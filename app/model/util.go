package model

import (
	"github.com/coopernurse/gorp"
	"github.com/mantishK/wormhole/storage"
)

func MysqlConnection() (*gorp.DbMap, error) {
	dbMap := storage.MysqlConnection()
	dbMap.AddTableWithName(Todo{}, "todos").SetKeys(true, "TodoID")

	//remove the following in production as it adds uneccessary complexity.
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		return nil, err
	}
	return dbMap, nil
}
