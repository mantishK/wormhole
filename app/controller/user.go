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

type User struct {
}

func (u *User) GetUser(w http.ResponseWriter, r *http.Request, filterData map[string]interface{}) {
	view := views.NewView(w)
	dbMap, _, params := Init(w, r)
	defer dbMap.Db.Close()
	user := model.User{}
	user_id := filterData["user_id"].(int)
	userIdString := params.Get("user_id")
	var err error
	if userIdString != "" {
		user_id, err = strconv.Atoi(userIdString)
		if err != nil {
			view.RenderErrorJson(apperror.NewNotNumericInputError("", err, "user_id"))
			return
		}
	}
	user.UserID = user_id
	err = user.Get(dbMap)
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}

	result := make(map[string]interface{})

	result["user"] = user
	view.RenderJson(result)
}

func (u *User) UserNameExists(w http.ResponseWriter, r *http.Request, filterData map[string]interface{}) {
	view := views.NewView(w)
	dbMap, _, params := Init(w, r)
	defer dbMap.Db.Close()
	user := model.User{}
	requiredFields := []string{"user_name"}
	count, err := validate.RequiredParams(params, requiredFields)
	if err != nil {
		view.RenderErrorJson(apperror.NewRequiredError(err.Error(), requiredFields[count]))
		return
	}
	userName := params.Get("user_name")
	if !validate.UserName(userName) {
		view.RenderErrorJson(apperror.NewInvalidInputError("Must not contain spaces", "user_name"))
		return
	}
	user.UserName = userName
	exists, err := user.UserNameExists(dbMap)
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}

	result := make(map[string]interface{})
	result["exists"] = exists
	view.RenderJson(result)
}

func (u *User) SaveUser(w http.ResponseWriter, r *http.Request, filterData map[string]interface{}) {
	view := views.NewView(w)
	dbMap, data, _ := Init(w, r)
	defer dbMap.Db.Close()
	user := model.User{}
	requiredFields := []string{"user_name", "password"}
	count, err := validate.RequiredData(data, requiredFields)
	if err != nil {
		view.RenderErrorJson(apperror.NewRequiredError(err.Error(), requiredFields[count]))
		return
	}
	user.UserName = data["user_name"].(string)
	user.Password = data["password"].(string)

	if !validate.UserName(user.UserName) {
		view.RenderErrorJson(apperror.NewInvalidInputError("Invalid Username", "user_name"))
		return
	}
	if !validate.Password(user.Password) {
		view.RenderErrorJson(apperror.NewInvalidInputError("Password not valid", "password"))
		return
	}
	exists, err := user.UserNameExists(dbMap)
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	if exists {
		view.RenderErrorJson(apperror.NewUserNameExistsError("", "user_name"))
		return
	}
	user.HashPassword("")
	err = user.Save(dbMap)
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	err = user.Get(dbMap)
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	result := make(map[string]interface{})
	result["user"] = user
	view.RenderJson(result)
}

func (u *User) SignIn(w http.ResponseWriter, r *http.Request, filterData map[string]interface{}) {
	view := views.NewView(w)
	dbMap, data, _ := Init(w, r)
	defer dbMap.Db.Close()
	user := model.User{}
	requiredFields := []string{"user_name", "password"}
	count, err := validate.RequiredData(data, requiredFields)
	if err != nil {
		view.RenderErrorJson(apperror.NewRequiredError(err.Error(), requiredFields[count]))
		return
	}
	user.UserName = data["user_name"].(string)
	err = user.GetUserFromUserName(dbMap)
	if err == sql.ErrNoRows {
		view.RenderErrorJson(apperror.NewInvalidUserNamePasswordError(""))
		return
	}
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	salt, err := user.GetPasswordSalt(dbMap)
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	user.Password = data["password"].(string)
	user.HashPassword(salt)
	exists, err := user.IsValidUser(dbMap)
	if err != nil || !exists {
		view.RenderErrorJson(apperror.NewInvalidUserNamePasswordError(""))
		return
	}

	userToken := model.UserToken{}
	userToken.UserId = user.UserID
	err = userToken.Add(dbMap)
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	result := make(map[string]interface{})
	result["user_token"] = userToken.Token
	result["user_id"] = userToken.UserId
	view.RenderJson(result)
}

func (u *User) SignOut(w http.ResponseWriter, r *http.Request, filterData map[string]interface{}) {
	view := views.NewView(w)
	dbMap, _, _ := Init(w, r)
	defer dbMap.Db.Close()
	userToken := model.UserToken{}
	userToken.Token = filterData["user_token"].(string)
	err := userToken.Delete(dbMap)
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	result := make(map[string]interface{})
	result["user_token"] = userToken.Token
	view.RenderJson(result)
}

func (u *User) RenewToken(w http.ResponseWriter, r *http.Request, filterData map[string]interface{}) {
	view := views.NewView(w)
	dbMap, _, _ := Init(w, r)
	defer dbMap.Db.Close()
	userToken := model.UserToken{}
	userToken.Token = filterData["user_token"].(string)
	userToken.GetUserIdFromToken(dbMap)
	err := userToken.Update(dbMap)
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	result := make(map[string]interface{})
	result["user_token"] = userToken.Token
	view.RenderJson(result)
}

func (u *User) UpdatePassword(w http.ResponseWriter, r *http.Request, filterData map[string]interface{}) {
	view := views.NewView(w)
	dbMap, data, _ := Init(w, r)
	defer dbMap.Db.Close()
	user := model.User{}
	requiredFields := []string{"password", "new_password"}
	count, err := validate.RequiredData(data, requiredFields)
	if err != nil {
		view.RenderErrorJson(apperror.NewRequiredError(err.Error(), requiredFields[count]))
		return
	}
	user.UserID = filterData["user_id"].(int)
	user.Password = data["password"].(string)
	password := data["new_password"].(string)
	if !validate.Password(password) {
		view.RenderErrorJson(apperror.NewInvalidInputError("Password not valid", "NewPassword"))
		return
	}

	salt, err := user.GetPasswordSalt(dbMap)
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	user.HashPassword(salt)

	exists, err := user.IsValidPassword(dbMap)
	if err != nil || !exists {
		view.RenderErrorJson(apperror.NewInvalidPasswordError("", "password"))
		return
	}

	user.Password = password
	user.HashPassword(salt)

	err = user.UpdatePassword(dbMap)
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	err = user.Get(dbMap)
	if err != nil {
		view.RenderErrorJson(apperror.NewDBError("", err))
		return
	}
	result := make(map[string]interface{})
	result["user"] = user
	view.RenderJson(result)
}
