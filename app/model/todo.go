package model

import (
	"database/sql"
	"log"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/jmoiron/sqlx"
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

func (t *Todo) SaveSqlx(dbMap *sqlx.DB) error {
	t.Created = time.Now()
	t.Modified = time.Now()
	_, err := dbMap.NamedExec("INSERT INTO todos (title, is_completed, created, modified) VALUES (:title, :is_completed, :created, :modified)", t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Todo) UpdateSqlx(dbMap *sqlx.DB) error {
	t.Modified = time.Now()
	_, err := dbMap.NamedExec("UPDATE todos set title = :title, is_completed = :is_completed, created = :created, modified = :modified WHERE todo_id = :todo_id", t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Todo) DeleteSqlx(dbMap *sqlx.DB) error {
	t.Created = time.Now()
	t.Modified = time.Now()
	_, err := dbMap.NamedExec("DELETE FROM todos WHERE todo_id = :todo_id", t)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllSqlx(dbMap *sqlx.DB) error {
	_, err := dbMap.Exec("DELETE FROM todos")
	if err != nil {
		return err
	}
	return nil
}

func (t *Todo) GetSqlx(dbMap *sqlx.DB) error {
	err := dbMap.Get(t, "SELECT * FROM todos WHERE todo_id = $1", t.TodoID)
	if err != nil {
		return err
	}
	return nil
}

func GetAllTodosSqlx(dbMap *sqlx.DB, offset, count int) ([]Todo, int, error) {
	var todos []Todo
	err := dbMap.Select(&todos, "SELECT SQL_CALC_FOUND_ROWS * FROM todos LIMIT ? , ? ", offset, count)
	if err != nil {
		return nil, 0, err
	}
	row := dbMap.QueryRowx("SELECT FOUND_ROWS()")
	var total int
	err = row.Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	return todos, int(total), nil
}

func (t *Todo) SaveSql(dbMap *sql.DB) error {
	t.Created = time.Now()
	t.Modified = time.Now()
	_, err := dbMap.Exec("INSERT INTO todos (title, is_completed, created, modified) VALUES (?, ?, ?, ?)", t.Title, t.IsCompleted, t.Created, t.Modified)
	if err != nil {
		return err
	}
	return nil
}

func (t *Todo) UpdateSql(dbMap *sql.DB) error {
	t.Modified = time.Now()
	_, err := dbMap.Exec("UPDATE todos set title = ?, is_completed = ?, created = ?, modified = ? WHERE todo_id = ?", t.Title, t.IsCompleted, t.Created, t.Modified, t.TodoID)
	if err != nil {
		return err
	}
	return nil
}

func (t *Todo) DeleteSql(dbMap *sql.DB) error {
	t.Created = time.Now()
	t.Modified = time.Now()
	_, err := dbMap.Exec("DELETE FROM todos WHERE todo_id = ?", t.TodoID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllSql(dbMap *sql.DB) error {
	_, err := dbMap.Exec("DELETE FROM todos")
	if err != nil {
		return err
	}
	return nil
}

func (t *Todo) GetSql(dbMap *sql.DB) error {
	err := dbMap.QueryRow("SELECT * FROM todos WHERE todo_id = ?", t.TodoID).Scan(t)
	if err != nil {
		return err
	}
	return nil
}

func GetAllTodosSql(dbMap *sql.DB, offset, count int) ([]Todo, int, error) {
	var todos []Todo
	todos = make([]Todo, 100)
	rows, err := dbMap.Query("SELECT SQL_CALC_FOUND_ROWS * FROM todos LIMIT ? , ? ", offset, count)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.TodoID, &todo.Title, &todo.IsCompleted, &todo.Created, &todo.Modified); err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}

	var total int
	err = dbMap.QueryRow("SELECT FOUND_ROWS()").Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	return todos, int(total), nil
}
