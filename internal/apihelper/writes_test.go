package apihelper

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func validateCodeAndHeader(resp *http.Response, t *testing.T, statusCode int, contentType string) {

	if statusCode != resp.StatusCode {
		t.Errorf("statusCode mismatch: expected %v, actual %v", statusCode, resp.StatusCode)
	}
	if contentType != resp.Header.Get("Content-Type") {
		t.Errorf("ContentType mismatch: expected %v, actual %v",
			contentType, resp.Header.Get("Content-Type"))
	}
}

func TestWriteJSON(t *testing.T) {
	expectedStatusCode := 200
	expectedContentType := "application/json; charset=UTF-8"
	expectedBody := "{\"Val\":2}\n"

	w := httptest.NewRecorder()
	data := struct{ Val int }{2}
	statusCode := 200
	WriteJSON(w, data, statusCode)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	validateCodeAndHeader(resp, t, expectedStatusCode, expectedContentType)

	if expectedBody != string(body) {
		t.Errorf("statusCode mismatch: expected [%v], actual [%v]",
			expectedBody, string(body))
	}
}

func TestWriteJSONFromText(t *testing.T) {
	expectedStatusCode := 200
	expectedContentType := "application/json; charset=UTF-8"
	expectedBody := "42"

	w := httptest.NewRecorder()
	data := []byte("42")
	statusCode := 200
	WriteJSONFromText(w, data, statusCode)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	validateCodeAndHeader(resp, t, expectedStatusCode, expectedContentType)

	if expectedBody != string(body) {
		t.Errorf("statusCode mismatch: expected [%v], actual [%v]",
			expectedBody, string(body))
	}
}

func TestWriteOK(t *testing.T) {
	expectedStatusCode := 200
	expectedContentType := "application/json; charset=UTF-8"
	expectedBody := "{\"Val\":2}\n"

	w := httptest.NewRecorder()
	data := struct{ Val int }{2}
	WriteOK(w, data)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	validateCodeAndHeader(resp, t, expectedStatusCode, expectedContentType)

	if expectedBody != string(body) {
		t.Errorf("statusCode mismatch: expected [%v], actual [%v]",
			expectedBody, string(body))
	}
}

func TestWriteError(t *testing.T) {
	expectedStatusCode := 500
	expectedContentType := "application/json; charset=UTF-8"
	expectedBody := `{"error":"message","status_code":500}
`

	w := httptest.NewRecorder()
	data := "message"
	statusCode := 500
	WriteError(w, data, statusCode)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	validateCodeAndHeader(resp, t, expectedStatusCode, expectedContentType)

	if expectedBody != string(body) {
		t.Errorf("statusCode mismatch: expected [%v], actual [%v]",
			expectedBody, string(body))
	}
}

func TestWrite404(t *testing.T) {
	expectedStatusCode := 404
	expectedContentType := "application/json; charset=UTF-8"
	expectedBody := `{"error":"item not found","status_code":404}
`

	w := httptest.NewRecorder()
	Write404(w)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	validateCodeAndHeader(resp, t, expectedStatusCode, expectedContentType)

	if expectedBody != string(body) {
		t.Errorf("statusCode mismatch: expected [%v], actual [%v]",
			expectedBody, string(body))
	}
}

func TestWrite500(t *testing.T) {
	expectedStatusCode := 500
	expectedContentType := "application/json; charset=UTF-8"
	expectedBody := `{"error":"message","status_code":500}
`

	w := httptest.NewRecorder()
	data := "message"
	Write500(w, errors.New(data))

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	validateCodeAndHeader(resp, t, expectedStatusCode, expectedContentType)

	if expectedBody != string(body) {
		t.Errorf("statusCode mismatch: expected [%v], actual [%v]",
			expectedBody, string(body))
	}
}
