package controllers

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

type TestCase struct {
	ID         string
	Response   string
	StatusCode int
}

func TestGet(t *testing.T) {
	for caseNum, item := range createCases() {
		url := "https://localhost:8080/post?id=" + item.ID
		//log.Printf("url %v", url)
		//req := httptest.NewRequest("GET", url, nil)
		//w := httptest.NewRecorder()

		//GetPost3(w, req)
		resp, err := http.Get(url)
		log.Printf("url %v", url)
		if err != nil {
			t.Errorf("error %v", err)
		}
		defer resp.Body.Close()

		//if w.Code != item.StatusCode {
		//	t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
		//		caseNum, w.Code, item.StatusCode)
		//}

		//resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}

		bodyStr := string(body)
		if bodyStr != item.Response {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				caseNum, bodyStr, item.Response)
		}
	}
}

func createCases() []TestCase {
	return []TestCase{
		{
			ID:         "3",
			Response:   `{"status": 200, "resp": {"post": 3}}`,
			StatusCode: http.StatusOK,
		},
		{
			ID:         "500",
			Response:   `{"status": 500, "err": "db_error"}`,
			StatusCode: http.StatusInternalServerError,
		},
	}
}
