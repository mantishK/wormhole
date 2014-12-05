package controller

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/mantishK/wormhole/app/apperror"
	"github.com/mantishK/wormhole/app/model"
	"github.com/mantishK/wormhole/app/validate"
	"github.com/mantishK/wormhole/app/views"
)

type Todo struct {
}

func (t *Todo) Add(w http.ResponseWriter, r *http.Request, filterData map[string]interface{}) {
	view := views.NewView(w)
	dbMap, data, _ := Init(w, r)
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

	todo.UserID = filterData["user_id"].(int)

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

func (t *Todo) Update(w http.ResponseWriter, r *http.Request, filterData map[string]interface{}) {
	view := views.NewView(w)
	dbMap, data, _ := Init(w, r)
	defer dbMap.Db.Close()
	var err error
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

	//check if the todo belongs to the user
	userID := filterData["user_id"].(int)
	if userID != todo.UserID {
		view.RenderHttpError("You are unauthorized!!", 401)
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

func (t *Todo) Delete(w http.ResponseWriter, r *http.Request, filterData map[string]interface{}) {
	view := views.NewView(w)
	dbMap, _, params := Init(w, r)
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

	//check if the todo belongs to the user
	userID := filterData["user_id"].(int)
	if userID != todo.UserID {
		view.RenderHttpError("You are unauthorized!!", 401)
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

func (t *Todo) Get(w http.ResponseWriter, r *http.Request, filterData map[string]interface{}) {
	view := views.NewView(w)
	dbMap, _, params := Init(w, r)
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

	//check if the todo belongs to the user
	userID := filterData["user_id"].(int)
	if userID != todo.UserID {
		view.RenderHttpError("You are unauthorized!!", 401)
	}

	result["todo"] = todo
	view.RenderJson(result)
}

func (t *Todo) GetAllUserTodos(w http.ResponseWriter, r *http.Request, filterData map[string]interface{}) {
	view := views.NewView(w)
	dbMap, _, params := Init(w, r)
	defer dbMap.Db.Close()

	offset := 0
	count := 10
	if params.Get("offset") != "" {
		offset, _ = strconv.Atoi(params.Get("offset"))
	}
	if params.Get("count") != "" {
		count, _ = strconv.Atoi(params.Get("count"))
	}

	userID := filterData["user_id"].(int)

	todos, total, err := model.GetAllUserTodos(dbMap, userID, offset, count)
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
