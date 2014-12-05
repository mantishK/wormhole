package filters

import (
	// "fmt"

	// "github.com/mantishK/wormhole/core/apperror"

	"net/http"
	"time"

	"github.com/mantishK/wormhole/app/apperror"
	"github.com/mantishK/wormhole/app/model"
	"github.com/mantishK/wormhole/app/views"
	"github.com/mantishK/wormhole/config"
)

type Authorize struct {
}

func (a *Authorize) Filter(writer http.ResponseWriter, req *http.Request, filterData map[string]interface{}) bool {
	token := req.Header.Get("X-TOKEN")
	view := views.NewView(writer)
	if token == "" {
		view.RenderHttpError("You are unauthorized!!", 401)
		return false
	}
	dbMap := config.NewConnection()
	defer dbMap.Db.Close()
	userToken := model.UserToken{}
	userToken.Token = token
	err := userToken.GetUserIdFromToken(dbMap)
	if err != nil || userToken.UserId == 0 {
		view.RenderHttpError("You are unauthorized!!", 401)
		return false
	}
	now := time.Now()
	timeDifference, err := time.ParseDuration("168h")
	if err != nil {
		view.RenderHttpError("This is not good, something went wrong!!", 500)
		return false
	}
	timeDifferencePlusModified := userToken.Modified.Add(timeDifference)
	if now.After(timeDifferencePlusModified) {
		view.RenderErrorJson(apperror.NewTokenExpiredError(""))
		return false
	}
	filterData["user_id"] = userToken.UserId
	filterData["user_token"] = token
	return true
}
