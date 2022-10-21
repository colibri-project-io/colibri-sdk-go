package webrest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const (
	timeoutDefault          uint16 = 1
	rest_client_transaction string = "Rest-Client"
)

type RestClient struct {
	name    string
	baseURL string
	client  *http.Client
}

func NewRestClient(name, baseURL string, timeout uint16) *RestClient {
	if timeout == 0 {
		timeout = timeoutDefault
	}

	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	client.Transport = newrelic.NewRoundTripper(client.Transport)
	return &RestClient{
		name:    name,
		baseURL: baseURL,
		client:  client,
	}
}

func Get[T interface{}](ctx context.Context, client *RestClient, path string, headers map[string]string) (*T, error) {
	return callRequestWithoutBody[T](ctx, http.MethodGet, client, path, headers)
}

func Post[T interface{}](ctx context.Context, client *RestClient, path string, body interface{}, headers map[string]string) (*T, error) {
	return callRequest[T](ctx, http.MethodPost, client, path, body, headers)
}

func Put[T interface{}](ctx context.Context, client *RestClient, path string, body interface{}, headers map[string]string) (*T, error) {
	return callRequest[T](ctx, http.MethodPut, client, path, body, headers)
}

func Delete[T interface{}](ctx context.Context, client *RestClient, path string, headers map[string]string) (*T, error) {
	return callRequestWithoutBody[T](ctx, http.MethodDelete, client, path, headers)
}

func callRequestWithoutBody[T interface{}](ctx context.Context, method string, client *RestClient, path string, headers map[string]string) (_ *T, err error) {
	txn := monitoring.GetTransactionInContext(ctx)
	if txn != nil {
		segment := monitoring.StartTransactionSegment(txn, rest_client_transaction, map[string]interface{}{
			"method": method,
			"path":   path,
		})
		defer monitoring.EndTransactionSegment(segment)
	}

	url := fmt.Sprintf("%s%s", client.baseURL, path)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create new http request %s: %v", url, err)
	}

	addHeaders(req, headers)
	resp, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return nil, fmt.Errorf("%d statusCode", resp.StatusCode)
	}

	return decodeResponse[T](resp)
}

func callRequest[T interface{}](ctx context.Context, method string, client *RestClient, path string, body interface{}, headers map[string]string) (_ *T, err error) {
	txn := monitoring.GetTransactionInContext(ctx)
	if txn != nil {
		segment := monitoring.StartTransactionSegment(txn, rest_client_transaction, map[string]interface{}{
			"method": method,
			"path":   path,
		})
		defer monitoring.EndTransactionSegment(segment)
	}

	url := fmt.Sprintf("%s%s", client.baseURL, path)
	bytesBody, err := makeBytesBody(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytesBody)
	if err != nil {
		return nil, fmt.Errorf("could not create new http request %s: %v", url, err)
	}

	addHeaders(req, headers)
	resp, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return nil, fmt.Errorf("%d statusCode", resp.StatusCode)
	}

	return decodeResponse[T](resp)
}

func makeBytesBody(body interface{}) (*bytes.Buffer, error) {
	if body == nil {
		return nil, nil
	}

	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("could not marshal entity body: %v", err)
	}

	return bytes.NewBuffer(requestBody), nil
}

func decodeResponse[T interface{}](resp *http.Response) (*T, error) {
	defer func() {
		closeBody(resp)
	}()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	var responseModel T
	if err := json.NewDecoder(resp.Body).Decode(&responseModel); err != nil {
		return nil, fmt.Errorf("could not decode response: %v", err)
	}
	return &responseModel, nil
}

func addHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func closeBody(resp *http.Response) {
	if resp == nil {
		return
	}

	if err := resp.Body.Close(); err != nil {
		logging.Error("error when close response body: %v\n", err)
	}
}
