package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Keshav-Agrawal/mongoapi/model"
)

func TestGetMyAllTask(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/tasks", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	rec := httptest.NewRecorder()
	GetMyAllTask(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("%v", res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	str := string(b)
	fmt.Println(str)
	if err != nil {
		t.Fatalf("%v", err)
	}
	//fmt.Println("b", b)
	//c := client.Ping(context.Background(), nil)
	//t.Log(c)
}
func TestInsertTask(t *testing.T) {
	/*arg := model.Homework{
		Task: "tom",
		Done: false,
	}*/
	postBody, _ := json.Marshal(model.Homework{
		Task: "tom",
		Done: false,
	})
	responseBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest("POST", "/api/task", responseBody)
	if err != nil {
		t.Fatalf("%v", err)
	}
	rec := httptest.NewRecorder()
	CreateTask(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("%v", res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	str := string(b)
	fmt.Println(str)
	if err != nil {
		t.Fatalf("%v", err)
	}
	//fmt.Println("b", b)
	//c := client.Ping(context.Background(), nil)
	//t.Log(c)
}

func TestUpdateTask(t *testing.T) {

	postBody, _ := json.Marshal(model.Homework{
		Task: "1234",
		Done: true,
	})
	client := &http.Client{}
	json, err := json.Marshal(postBody)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPut, "http://localhost:4000/api/task/621f1031ddc5de508f36d0db", bytes.NewBuffer(json))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.StatusCode)
	/*MarkAsDone(resp, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("%v", res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	str := string(b)
	fmt.Println(str)
	if err != nil {
		t.Fatalf("%v", err)
	}
	//fmt.Println("b", b)
	//c := client.Ping(context.Background(), nil)
	//t.Log(c)
	*/
}


