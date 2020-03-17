package controllers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	ID         string
	Response   string
	StatusCode int
}

func TestGet(t *testing.T) {
	for caseNum, item := range createCases() {
		url := "http://localhost/post?id=" + item.ID
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		Get(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result()
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
			ID:         "42",
			Response:   `{"status": 200, "resp": {"user": 42}}`,
			StatusCode: http.StatusOK,
		},
		{
			ID:         "500",
			Response:   `{"status": 500, "err": "db_error"}`,
			StatusCode: http.StatusInternalServerError,
		},
	}
}
