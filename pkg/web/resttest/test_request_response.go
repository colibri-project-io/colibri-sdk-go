package resttest

import (
	"encoding/json"
	"io"
	"net/http/httptest"
)

// RequestTestResponse is a contract to test http requests responses
type RequestTestResponse struct {
	writer *httptest.ResponseRecorder
}

// StatusCode returns the result status code
func (r *RequestTestResponse) StatusCode() int {
	return r.writer.Result().StatusCode
}

// StatusCode returns the result body decoded
func (r *RequestTestResponse) DecodeBody(object any) error {
	body, err := io.ReadAll(r.writer.Result().Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, object)
}
