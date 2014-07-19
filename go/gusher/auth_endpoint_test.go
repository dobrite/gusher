package gusher

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	resp, req := getRespReq(t, "GET", "http://localhost:3000/pusher/auth", nil)
	NewServeMux("app", "tester").ServeHTTP(resp, req)
	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fail()
	} else {
		if resp.Code != 405 {
			t.Error("Status code not 405: Method not allowed")
		}
	}
}

func TestAuthEndpointJSONP(t *testing.T) {
	data := url.Values{}
	data.Set("callback", "func")
	data.Set("socket_id", "1234.1234")
	data.Set("channel_name", "private-foobar")
	data.Set("channel_data", "yo")
	encoded := data.Encode()

	postBody := bytes.NewBufferString(encoded)

	resp, req := getRespReq(t, "POST", "http://localhost:3000/pusher/auth", postBody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	NewServeMux("app", "tester").ServeHTTP(resp, req)
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fail()
	} else {
		log.Println(string(body))
	}
}
