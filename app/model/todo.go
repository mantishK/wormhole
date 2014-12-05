package model

import (
	"time"

	"github.com/coopernurse/gorp"
)

type Todo struct {
	TodoID      int       `db:"todo_id" json:"todo_id"`
	UserID      int       `db:"user_id" json:"user_id"`
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

func (t *Todo) Get(dbMap *gorp.DbMap) error {
	err := dbMap.SelectOne(t, "SELECT * FROM todos WHERE todo_id = ?", t.TodoID)
	if err != nil {
		return err
	}
	return nil
}

func GetAllUserTodos(dbMap *gorp.DbMap, userID, offset, count int) ([]Todo, int, error) {
	var todos []Todo
	_, err := dbMap.Select(&todos, "SELECT SQL_CALC_FOUND_ROWS * FROM todos WHERE user_id = ? LIMIT ?,?", userID, offset, count)
	if err != nil {
		return nil, 0, err
	}
	total, err := dbMap.SelectInt("SELECT FOUND_ROWS()")
	if err != nil {
		return nil, 0, err
	}
	return todos, int(total), nil
}
