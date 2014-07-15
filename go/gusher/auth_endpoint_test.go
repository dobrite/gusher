package gusher

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getRespReq(t *testing.T, method string, url string, body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	resp := httptest.NewRecorder()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	return resp, req
}

func TestAuthEndpointGET(t *testing.T) {
	resp, req := getRespReq(t, "POST", "http://localhost:3000/pusher/auth", nil)
	NewServeMux("app", "tester").ServeHTTP(resp, req)
	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fail()
	} else {
		if resp.Code != 405 {
			t.Error("Status code not 405: Method not allowed")
		}
	}
}

//func TestAuthEndpointJSONP(t *testing.T) {
//	resp, req := getRespReq(t, "POST", "http://localhost:3000/pusher/auth", nil)
//	NewServeMux("app", "tester").ServeHTTP(resp, req)
//	if _, err := ioutil.ReadAll(resp.Body); err != nil {
//		t.Fail()
//	} else {
//		if resp.Code != 405 {
//			t.Error("Status code not 405: Method not allowed")
//		}
//	}
//}
