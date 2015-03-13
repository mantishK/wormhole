package controller

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/mantishK/httprouter"
	"github.com/mantishK/wormhole/app/apperror"
	"github.com/mantishK/wormhole/app/model"
	"github.com/mantishK/wormhole/app/validate"
	"github.com/mantishK/wormhole/app/views"
)

type Todo struct {
}

func (t *Todo) Add(w http.ResponseWriter, r *http.Request, _ httprouter.Params, filterData map[string]interface{}) {
	// func (t *Todo) Add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	view := views.NewView(w)
	data, _ := DataParse(w, r)
	dbMap, err := model.MysqlConnection()

	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	defer dbMap.Db.Close()

	todo := model.Todo{}
	requiredFields := []string{"title"}
	count, err := validate.RequiredData(data, requiredFields)
	if err != nil {
		view.RenderErrorJson(apperror.NewRequiredError(err.Error(), requiredFields[count]))
		return
	}
	//required field assignment
	todo.Title = data["title"].(string)

	//optional field assignment
	// if data["something"] != nil {
	//   todo.something = int(data["something"].(float64))
	// }

	//validate
	if data["title"] == nil {
		view.RenderErrorJson(apperror.NewInvalidInputError("Empty Title", "title"))
		return
	}

	//save todo
	err = todo.Save(dbMap)
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}

	//build and return json
	result := make(map[string]interface{})
	result["todo"] = todo
	view.RenderJson(result)
}

func (t *Todo) Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params, filterData map[string]interface{}) {
	// func (t *Todo) Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	view := views.NewView(w)
	data, _ := DataParse(w, r)
	dbMap, err := model.MysqlConnection()
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	defer dbMap.Db.Close()
	requiredFields := []string{"todo_id", "title", "isCompleted"}
	count, err := validate.RequiredData(data, requiredFields)
	if err != nil {
		view.RenderErrorJson(apperror.NewRequiredError(err.Error(), requiredFields[count]))
		return
	}
	todo := model.Todo{}
	todo.TodoID = int(data["todo_id"].(float64))
	err = todo.Get(dbMap)
	if err == sql.ErrNoRows {
		view.RenderHttpError("Todo not found.", http.StatusNotFound)
		return
	}
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	todo.Title = data["title"].(string)
	todo.IsCompleted = data["isCompleted"].(bool)
	err = todo.Update(dbMap)
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}

	result := make(map[string]interface{})
	result["todo"] = todo
	view.RenderJson(result)
}

func (t *Todo) Delete(w http.ResponseWriter, r *http.Request, _ httprouter.Params, filterData map[string]interface{}) {
	// func (t *Todo) Delete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	view := views.NewView(w)
	_, params := DataParse(w, r)
	dbMap, err := model.MysqlConnection()
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	defer dbMap.Db.Close()
	requiredFields := []string{"id"}
	count, err := validate.RequiredParams(params, requiredFields)
	if err != nil {
		view.RenderErrorJson(apperror.NewRequiredError(err.Error(), requiredFields[count]))
		return
	}
	todo := model.Todo{}
	todoIdString := params.Get("id")
	if todoIdString != "" {
		todoID, err := strconv.Atoi(todoIdString)
		if err != nil {
			view.RenderErrorJson(apperror.NewNotNumericInputError("", err, "id"))
			return
		}
		todo.TodoID = todoID
	}

	err = todo.Get(dbMap)
	if err == sql.ErrNoRows {
		view.RenderHttpError("Todo not found.", http.StatusNotFound)
		return
	}
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	todo.Delete(dbMap)
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	result := make(map[string]interface{})

	result["todo"] = todo
	view.RenderJson(result)
}

func (t *Todo) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params, filterData map[string]interface{}) {
	// func (t *Todo) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	view := views.NewView(w)
	_, params := DataParse(w, r)
	dbMap, err := model.MysqlConnection()
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	defer dbMap.Db.Close()
	requiredFields := []string{"todo_id"}
	count, err := validate.RequiredParams(params, requiredFields)
	if err != nil {
		view.RenderErrorJson(apperror.NewRequiredError(err.Error(), requiredFields[count]))
		return
	}
	todo := model.Todo{}
	todoIdString := params.Get("todo_id")
	if todoIdString != "" {
		todoId, err := strconv.Atoi(todoIdString)
		if err != nil {
			view.RenderErrorJson(apperror.NewNotNumericInputError("", err, "user_id"))
			return
		}
		todo.TodoID = todoId
	}
	err = todo.Get(dbMap)
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}

	result := make(map[string]interface{})

	result["todo"] = todo
	view.RenderJson(result)
}

func (t *Todo) GetAllTodos(w http.ResponseWriter, r *http.Request, _ httprouter.Params, filterData map[string]interface{}) {
	// func (t *Todo) GetAllTodos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	view := views.NewView(w)
	_, params := DataParse(w, r)
	dbMap, err := model.MysqlConnection()
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	defer dbMap.Db.Close()

	offset := 0
	count := 10
	if params.Get("offset") != "" {
		offset, _ = strconv.Atoi(params.Get("offset"))
	}
	if params.Get("count") != "" {
		count, _ = strconv.Atoi(params.Get("count"))
	}

	todos, total, err := model.GetAllTodos(dbMap, offset, count)
	if err == sql.ErrNoRows {
		view.RenderHttpError("No todos found.", http.StatusNotFound)
		return
	} else if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}

	result := make(map[string]interface{})
	meta := make(map[string]interface{})
	result["todos"] = todos
	meta["total"] = total
	result["meta"] = meta
	view.RenderJson(result)
}
