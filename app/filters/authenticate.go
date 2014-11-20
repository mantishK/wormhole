package filters

import (
	"net/http"

	"github.com/mantishK/wormhole/app/views"
)

type Authenticate struct {
}

func (a *Authenticate) Filter(writer http.ResponseWriter, req *http.Request, filterData map[string]interface{}) bool {
	//Todo: implement the api authentication

	Authenticated := true
	view := views.NewView(writer)
	if !Authenticated {
		view.RenderHttpError("You are authenticated!!", 401)
		return false
	}
	return Authenticated
}
