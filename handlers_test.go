package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	timing_data = `{"page-uri":"/foo/bar",
		"nav-timing": {
			"dns":1,
			"connect":2,
			"ttfb":3,
			"basePage":4,
			"frontEnd":5
		}
	}`
	js_data = `{"page-uri": "fizz/buzz",
		"query-string": "param=value&other=not",
		"js-error": {
			"error-type": "ReferenceError",
			"description": "func is not defined"
		}
	}`
	csp_data = `{"csp-report": {
		"document-uri": "https://www.example.com/",
		"blocked-uri": "https://evil.example.com/",
		"violated-directive": "directive",
		"original-policy": "policy"
	}}`
	recorders = []Recorder{StatsDRecorder{}}
)

func TestNavTimingHandlerSuccess(t *testing.T) {
	req, _ := http.NewRequest("POST", "/r", bytes.NewBufferString(timing_data))
	req.Header.Add("X-Real-Ip", "192.168.0.1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:36.0) Gecko/20100101 Firefox/36.0")
	resp := httptest.NewRecorder()

	NavTimingHandler(recorders)(resp, req)

	const expected_response_code = 200
	if code := resp.Code; code != expected_response_code {
		t.Errorf("received %v response code, expected %v", code, expected_response_code)
	}
}

func TestNavTimingHandlerNotPOST(t *testing.T) {
	req, _ := http.NewRequest("GET", "/r", bytes.NewBufferString(timing_data))
	resp := httptest.NewRecorder()

	NavTimingHandler(recorders)(resp, req)

	const expected_response_code = 405
	if code := resp.Code; code != expected_response_code {
		t.Errorf("received %v response code, expected %v", code, expected_response_code)
	}
}

func TestNavTimingHandlerInvalidJSON(t *testing.T) {
	req, _ := http.NewRequest("POST", "/r", bytes.NewBufferString(`{invalid:"json"}`))
	resp := httptest.NewRecorder()

	NavTimingHandler(recorders)(resp, req)

	const expected_response_code = 400
	if code := resp.Code; code != expected_response_code {
		t.Errorf("expected %v, but received %v response code", expected_response_code, code)
	}
	const expected_error_message = "Error parsing JSON\n"
	if msg := resp.Body.String(); msg != expected_error_message {
		t.Errorf("expected \"%v\", but found \"%v\" error message", expected_error_message, msg)
	}
}

func TestNavTimingHandlerInvalidPageURI(t *testing.T) {
	req, _ := http.NewRequest("POST", "/r", bytes.NewBufferString(`{"page-uri":"/foo/bar///"}`))
	resp := httptest.NewRecorder()

	NavTimingHandler(recorders)(resp, req)

	const expected_response_code = 406
	if code := resp.Code; code != expected_response_code {
		t.Errorf("expected %v, but received %v response code", expected_response_code, code)
	}
	const expected_error_message = "Invalid page-uri passed\n"
	if msg := resp.Body.String(); msg != expected_error_message {
		t.Errorf("expected \"%v\", but found \"%v\" error message", expected_error_message, msg)
	}
}

func TestJsErrorHandlerSuccess(t *testing.T) {
	req, _ := http.NewRequest("POST", "/r", bytes.NewBufferString(js_data))
	req.Header.Add("X-Real-Ip", "192.168.0.1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:36.0) Gecko/20100101 Firefox/36.0")
	resp := httptest.NewRecorder()

	JsErrorReportHandler()(resp, req)

	const expected_response_code = 200
	if code := resp.Code; code != expected_response_code {
		t.Errorf("received %v response code, expected %v", code, expected_response_code)
	}
}

func TestJsErrorHandlerNotPOST(t *testing.T) {
	req, _ := http.NewRequest("GET", "/r", bytes.NewBufferString(js_data))
	resp := httptest.NewRecorder()

	JsErrorReportHandler()(resp, req)

	const expected_response_code = 405
	if code := resp.Code; code != expected_response_code {
		t.Errorf("received %v response code, expected %v", code, expected_response_code)
	}
}

func TestJsErrorHandlerInvalidJSON(t *testing.T) {
	req, _ := http.NewRequest("POST", "/r", bytes.NewBufferString(`{invalid:"json"}`))
	resp := httptest.NewRecorder()

	JsErrorReportHandler()(resp, req)

	const expected_response_code = 400
	if code := resp.Code; code != expected_response_code {
		t.Errorf("expected %v, but received %v response code", expected_response_code, code)
	}
	const expected_error_message = "Error parsing JSON\n"
	if msg := resp.Body.String(); msg != expected_error_message {
		t.Errorf("expected \"%v\", but found \"%v\" error message", expected_error_message, msg)
	}
}

func TestCSPReportHandlerSuccess(t *testing.T) {
	req, _ := http.NewRequest("POST", "/r", bytes.NewBufferString(csp_data))
	req.Header.Add("X-Real-Ip", "192.168.0.1")
	resp := httptest.NewRecorder()

	CSPReportHandler()(resp, req)

	const expected_response_code = 200
	if code := resp.Code; code != expected_response_code {
		t.Errorf("received %v response code, expected %v", code, expected_response_code)
	}
}

func TestCSPReportHandlerNotPOST(t *testing.T) {
	req, _ := http.NewRequest("GET", "/r", bytes.NewBufferString(csp_data))
	resp := httptest.NewRecorder()

	CSPReportHandler()(resp, req)

	const expected_response_code = 405
	if code := resp.Code; code != expected_response_code {
		t.Errorf("received %v response code, expected %v", code, expected_response_code)
	}
}

func TestCSPReportHandlerInvalidJSON(t *testing.T) {
	req, _ := http.NewRequest("POST", "/r", bytes.NewBufferString(`{invalid:"json"}`))
	resp := httptest.NewRecorder()

	CSPReportHandler()(resp, req)

	const expected_response_code = 400
	if code := resp.Code; code != expected_response_code {
		t.Errorf("expected %v, but received %v response code", expected_response_code, code)
	}
	const expected_error_message = "Error parsing JSON\n"
	if msg := resp.Body.String(); msg != expected_error_message {
		t.Errorf("expected \"%v\", but found \"%v\" error message", expected_error_message, msg)
	}
}
