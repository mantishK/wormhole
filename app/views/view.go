package views

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mantishK/wormhole/app/apperror"
)

type view struct {
	writer http.ResponseWriter
}

func NewView(w http.ResponseWriter) view {
	return view{w}
}

func (v *view) SetHeader(key, value string) {
	v.writer.Header().Set(key, value)
}

func (v *view) RenderJson(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	v.writer.Header().Set("Content-Type", "application/json")

	//remove in production
	v.writer.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Fprint(v.writer, string(jsonData))
	return nil
}

func (v *view) RenderErrorJson(apperror apperror.Apperror) error {
	result := make(map[string]interface{})
	result["id"] = apperror.GetIdString()
	result["message"] = apperror.GetMessage()
	if len(apperror.GetSysMesage()) > 0 {
		result["sys_message"] = apperror.GetSysMesage()
	} else {
		result["sys_message"] = nil
	}
	if len(apperror.GetField()) > 0 {
		result["field"] = apperror.GetField()
	} else {
		result["field"] = nil
	}
	result["response"] = "error"
	jsonData, err := json.Marshal(result)
	if err != nil {
		return err
	}
	// v.writer.WriteHeader(200)
	v.writer.Header().Set("Content-Type", "application/json")

	//remove in production
	v.writer.Header().Set("Access-Control-Allow-Origin", "*")

	v.writer.WriteHeader(apperror.GetHttpStatusCode())
	fmt.Fprint(v.writer, string(jsonData))
	return nil
}

func (v *view) RenderHttpError(message string, code int) error {
	v.writer.Header().Set("Access-Control-Allow-Origin", "*")
	http.Error(v.writer, message, code)
	return nil
}
