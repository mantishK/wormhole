package filters

import
// "fmt"

// "github.com/mantishK/wormhole/core/apperror"
"net/http"

type Authorize struct {
}

func (a *Authorize) Filter(writer http.ResponseWriter, req *http.Request, filterData map[string]interface{}) bool {
	//Todo: implement user authorization

	//You can send the filter data to the controllers using filterData
	//e.g filterData["user_id"] = getUserIdFromToken(req.Header.Get("X-TOKEN"))
	return true
}
