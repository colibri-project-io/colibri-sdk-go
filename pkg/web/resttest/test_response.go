package resttest

import (
	"encoding/json"
	"io"
	"net/http/httptest"
)

// TestResponse is a contract to test http requests responses
type TestResponse struct {
	writer *httptest.ResponseRecorder
}

// StatusCode returns the result status code
func (r *TestResponse) StatusCode() int {
	return r.writer.Result().StatusCode
}

// DecodeBody returns the result body decoded
func (r *TestResponse) DecodeBody(object any) error {
	body, err := io.ReadAll(r.writer.Result().Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, object)
}
