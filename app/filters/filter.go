package filters

import (
	"net/http"
)

type Filterable interface {
	Filter(http.ResponseWriter, *http.Request, map[string]interface{}) bool
}
