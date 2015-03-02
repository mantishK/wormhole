package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/mantishK/wormhole/app/views"
	"github.com/mantishK/wormhole/storage"
)

func Init(w http.ResponseWriter, r *http.Request) (map[string]interface{}, url.Values) {
	var data interface{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	json.Unmarshal(buf.Bytes(), &data)
	if data == nil {
		return nil, r.URL.Query()
	}
	return data.(map[string]interface{}), r.URL.Query()
}

func Terminate() {
	storage.CloseMysql()
}

func HasSpaces(token string) bool {
	return strings.Contains(token, " ")
}

func CorsOptionController(w http.ResponseWriter, r *http.Request, filterData map[string]interface{}) {
	v := views.NewView(w)
	//remove in production
	v.SetHeader("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE")
	v.SetHeader("Allow", "OPTIONS, GET, POST, PUT, DELETE")
	v.SetHeader("Access-Control-Allow-Headers", "accept, content-type")

	res := make(map[string]interface{})
	v.RenderJson(res)
}
