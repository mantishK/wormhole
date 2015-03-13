package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

var url = "http://localhost:8080"
var namespace = "/api"
var todoId int

func post() error {
	todo := make(map[string]interface{})
	todo["title"] = "Abc"
	resp, err := testCallApi(url+namespace+"/todos", "POST", todo)
	if err != nil {
		return err
	}
	todoResp := resp["todo"].(map[string]interface{})
	if todoResp["title"] != "Abc" {
		return errors.New("Todo not saved")
	}
	todoId = int(todoResp["todo_id"].(float64))
	return nil
}

func put() error {
	todo := make(map[string]interface{})
	todo["title"] = "Abc - 2"
	todo["todo_id"] = todoId
	todo["isCompleted"] = true
	resp, err := testCallApi(url+namespace+"/todos", "PUT", todo)
	if err != nil {
		return err
	}
	todoResp := resp["todo"].(map[string]interface{})
	if todoResp["title"] != "Abc - 2" {
		return errors.New("Todo not update")
	}
	return nil
}

func getAll() error {
	resp, err := testCallApi(url+namespace+"/todos?offset=0&count=10000", "GET", nil)
	if err != nil {
		return err
	}
	todoResp := resp["meta"].(map[string]interface{})
	if todoResp["total"] == 0 {
		return errors.New("No todos")
	}
	return nil
}

func get() error {
	resp, err := testCallApi(url+namespace+"/todo?todo_id="+strconv.Itoa(todoId), "GET", nil)
	if err != nil {
		return err
	}
	todoResp := resp["todo"].(map[string]interface{})
	if todoResp["title"] != "Abc - 2" && todoResp["isCompleted"] != true {
		return errors.New("Todo not update")
	}
	return nil
}

func delete() error {
	resp, err := testCallApi(url+namespace+"/todos?id="+strconv.Itoa(todoId), "DELETE", nil)
	if err != nil {
		return err
	}
	todoResp := resp["todo"].(map[string]interface{})
	if todoResp["title"] != "Abc - 2" && todoResp["isCompleted"] != true {
		return errors.New("Todo not update")
	}
	return nil
}

func BenchmarkPost(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := post()
		if err != nil {
			b.Error(err)
		}
	}
}

func TestPost(t *testing.T) {
	err := post()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkPut(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := put()
		if err != nil {
			b.Error(err)
		}
	}
}

func TestPut(t *testing.T) {
	err := put()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkGetAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := getAll()
		if err != nil {
			b.Error(err)
		}
	}
}

func TestGetAll(t *testing.T) {
	err := getAll()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := get()
		if err != nil {
			b.Error(err)
		}
	}
}

func TestGet(t *testing.T) {
	err := get()
	if err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	err := delete()
	if err != nil {
		t.Error(err)
	}
}

func testCallApi(u string, m string, d map[string]interface{}) (map[string]interface{}, error) {
	dByteArray, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(dByteArray)
	req, err := http.NewRequest(m, u, buf)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var bodyMap interface{}
	json.Unmarshal([]byte(body), &bodyMap)
	return bodyMap.(map[string]interface{}), nil
}
