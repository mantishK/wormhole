package model

import (
	"time"

	"github.com/coopernurse/gorp"
)

type Todo struct {
	TodoID      int       `db:"todo_id" json:"todo_id"`
	Title       string    `db:"title" json:"title"`
	IsCompleted bool      `db:"is_completed" json:"isCompleted"`
	Created     time.Time `db:"created" json:"created"`
	Modified    time.Time `db:"modified" json:"modified"`
}

func (t *Todo) Save(dbMap *gorp.DbMap) error {
	t.Created = time.Now()
	t.Modified = time.Now()
	err := dbMap.Insert(t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Todo) Update(dbMap *gorp.DbMap) error {
	t.Modified = time.Now()
	_, err := dbMap.Update(t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Todo) Delete(dbMap *gorp.DbMap) error {
	t.Created = time.Now()
	t.Modified = time.Now()
	_, err := dbMap.Delete(t)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAll(dbMap *gorp.DbMap) error {
	_, err := dbMap.Exec("DELETE FROM todos")
	if err != nil {
		return err
	}
	return nil
}

func (t *Todo) Get(dbMap *gorp.DbMap) error {
	err := dbMap.SelectOne(t, "SELECT * FROM todos WHERE todo_id = ?", t.TodoID)
	if err != nil {
		return err
	}
	return nil
}

func GetAllTodos(dbMap *gorp.DbMap, offset, count int) ([]Todo, int, error) {
	var todos []Todo
	_, err := dbMap.Select(&todos, "SELECT SQL_CALC_FOUND_ROWS * FROM todos LIMIT ?,?", offset, count)
	if err != nil {
		return nil, 0, err
	}
	total, err := dbMap.SelectInt("SELECT FOUND_ROWS()")
	if err != nil {
		return nil, 0, err
	}
	return todos, int(total), nil
}
