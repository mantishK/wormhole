package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/coopernurse/gorp"
	"github.com/mantishK/wormhole/app/views"
	"github.com/mantishK/wormhole/config"
)

func Init(w http.ResponseWriter, r *http.Request) (*gorp.DbMap, map[string]interface{}, url.Values) {
	dbMap := config.NewConnection()
	var data interface{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	json.Unmarshal(buf.Bytes(), &data)
	if data == nil {
		return dbMap, nil, r.URL.Query()
	}
	return dbMap, data.(map[string]interface{}), r.URL.Query()
}

func HasSpaces(token string) bool {
	return strings.Contains(token, " ")
}

func CorsOptionController(w http.ResponseWriter, r *http.Request, filterData map[string]interface{}) {
	view := views.NewView(w)
	//remove in production
	view.SetHeader("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE")
	view.SetHeader("Allow", "OPTIONS, GET, POST, PUT, DELETE")
	view.SetHeader("Access-Control-Allow-Headers", "accept, content-type")

	result := make(map[string]interface{})
	view.RenderJson(result)
}
