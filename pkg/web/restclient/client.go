package restclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
	"github.com/mercari/go-circuitbreaker"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const (
	timeoutDefault        uint   = 1
	restClientTransaction string = "REST-CLIENT"
)

var ErrServiceNotAvailable = errors.New("service not available")

type ResponseBody interface {
	any
}

type RequestData interface {
	any
}

type ResponseErrorData interface {
	any
}

// RestClient client http to encapsulate rest calls
type RestClient struct {
	name    string
	baseURL string
	client  *http.Client
	cb      *circuitbreaker.CircuitBreaker
}

// NewRestClient create new client http with timeout configuration and New Relic transport configuration
func NewRestClient(name string, baseURL string, timeout uint) *RestClient {
	if timeout == 0 {
		timeout = timeoutDefault
	}
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	client.Transport = newrelic.NewRoundTripper(client.Transport)
	return &RestClient{
		name:    name,
		baseURL: baseURL,
		client:  client,
		cb: circuitbreaker.New(
			circuitbreaker.WithOpenTimeout(time.Second*10),
			circuitbreaker.WithTripFunc(circuitbreaker.NewTripFuncConsecutiveFailures(5)),
			circuitbreaker.WithOnStateChangeHookFn(func(oldState, newState circuitbreaker.State) {
				logging.Info("[%s] state changed: old [%s] -> new [%s]", name, string(oldState), string(newState))
			}),
		),
	}
}

// Get call get request
func Get[T ResponseBody](ctx context.Context, client *RestClient, path string, headers map[string]string) ResponseData[T] {
	resp, _ := callRequestWithoutBody[T, any](ctx, http.MethodGet, client, path, headers)
	return resp
}

// GetWithErrorData call get request and when error occurs decode err into ResponseErrorData
func GetWithErrorData[T ResponseBody, E ResponseErrorData](ctx context.Context, client *RestClient, path string, headers map[string]string) (ResponseData[T], *E) {
	return callRequestWithoutBody[T, E](ctx, http.MethodGet, client, path, headers)
}

// Post call post request
func Post[T ResponseBody, R RequestData](ctx context.Context, client *RestClient, path string, entityBody *R, headers map[string]string) ResponseData[T] {
	resp, _ := callRequest[T, any](ctx, http.MethodPost, client, path, entityBody, headers)
	return resp
}

// PostWithErrorData call post request and when error occurs decode err into ResponseErrorData
func PostWithErrorData[T ResponseBody, E ResponseErrorData, R RequestData](ctx context.Context, client *RestClient, path string, entityBody *R, headers map[string]string) (ResponseData[T], *E) {
	return callRequest[T, E](ctx, http.MethodPost, client, path, entityBody, headers)
}

// PostBodyString call post request with string as body
// to set Content-type pass is header map
//
//	Ex.: to make a form url encoded request
//		`PostBodyString[MyReturnData](ctx, client, "/my-path", "body string for urlencoded...", map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
func PostBodyString[T ResponseBody](ctx context.Context, client *RestClient, path string, bodyString string, headers map[string]string) ResponseData[T] {
	return callRequestBodyString[T](ctx, http.MethodPost, client, path, bodyString, headers)
}

// Put call put request
func Put[T ResponseBody, R RequestData](ctx context.Context, client *RestClient, path string, entityBody *R, headers map[string]string) ResponseData[T] {
	resp, _ := callRequest[T, any](ctx, http.MethodPut, client, path, entityBody, headers)
	return resp
}

// PutWithErrorData call put request and when error occurs decode err into ResponseErrorData
func PutWithErrorData[T ResponseBody, E ResponseErrorData, R RequestData](ctx context.Context, client *RestClient, path string, entityBody *R, headers map[string]string) (ResponseData[T], *E) {
	return callRequest[T, E](ctx, http.MethodPut, client, path, entityBody, headers)
}

// Patch call patch request
func Patch[T ResponseBody, R RequestData](ctx context.Context, client *RestClient, path string, entityBody *R, headers map[string]string) ResponseData[T] {
	resp, _ := callRequest[T, any](ctx, http.MethodPatch, client, path, entityBody, headers)
	return resp
}

// PatchWithErrorData call patch request and when error occurs decode err into ResponseErrorData
func PatchWithErrorData[T ResponseBody, E ResponseErrorData, R RequestData](ctx context.Context, client *RestClient, path string, entityBody *R, headers map[string]string) (ResponseData[T], *E) {
	return callRequest[T, E](ctx, http.MethodPatch, client, path, entityBody, headers)
}

// Delete call delete request
func Delete[T ResponseBody](ctx context.Context, client *RestClient, path string, headers map[string]string) ResponseData[T] {
	resp, _ := callRequestWithoutBody[T, any](ctx, http.MethodDelete, client, path, headers)
	return resp
}

// DeleteWithErrorData call delete request and when error occurs decode err into ResponseErrorData
func DeleteWithErrorData[T ResponseBody, E ResponseErrorData](ctx context.Context, client *RestClient, path string, headers map[string]string) (ResponseData[T], *E) {
	return callRequestWithoutBody[T, E](ctx, http.MethodDelete, client, path, headers)
}

func callRequestWithoutBody[T ResponseBody, E ResponseErrorData](ctx context.Context, method string, client *RestClient, path string, headers map[string]string) (_ ResponseData[T], _ *E) {
	if !client.cb.Ready() {
		return newResponseData[T](http.StatusInternalServerError, nil, nil, ErrServiceNotAvailable), nil
	}

	var err error
	defer func() {
		err = client.cb.Done(ctx, err)
	}()

	txn := monitoring.GetTransactionInContext(ctx)
	if txn != nil {
		segment := createSegment(ctx, txn, method, path)
		defer monitoring.EndTransactionSegment(segment)
	}

	url := fmt.Sprintf("%s%s", client.baseURL, path)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fErr := fmt.Errorf("could not create new http request %s: %w", url, err)
		return newResponseData[T](http.StatusInternalServerError, nil, nil, fErr), nil
	}

	addHeaders(req, headers)
	resp, err := client.client.Do(req)
	if err != nil {
		return newResponseData[T](http.StatusInternalServerError, nil, nil, err), nil
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		respBErr, err := processErrorResponse[E](resp)
		return newResponseData[T](resp.StatusCode, nil, resp.Header, err), respBErr
	}

	r, err := decodeResponse[T](resp)
	return newResponseData(resp.StatusCode, r, resp.Header, err), nil
}

func callRequest[T ResponseBody, E ResponseErrorData, R RequestData](ctx context.Context, method string, client *RestClient, path string, entityBody *R, headers map[string]string) (_ ResponseData[T], _ *E) {
	if !client.cb.Ready() {
		return newResponseData[T](http.StatusInternalServerError, nil, nil, ErrServiceNotAvailable), nil
	}

	var err error
	defer func() {
		err = client.cb.Done(ctx, err)
	}()

	txn := monitoring.GetTransactionInContext(ctx)
	if txn != nil {
		segment := createSegment(ctx, txn, method, path)
		defer monitoring.EndTransactionSegment(segment)
	}

	url := fmt.Sprintf("%s%s", client.baseURL, path)
	bytesBody, err := makeBytesBody(entityBody)
	if err != nil {
		return newResponseData[T](http.StatusInternalServerError, nil, nil, err), nil
	}
	req, err := http.NewRequest(method, url, bytesBody)

	addHeaders(req, headers)
	resp, err := client.client.Do(req)
	if err != nil {
		return newResponseData[T](http.StatusInternalServerError, nil, nil, err), nil
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		errBody, err := processErrorResponse[E](resp)
		return newResponseData[T](resp.StatusCode, nil, resp.Header, err), errBody
	}

	r, err := decodeResponse[T](resp)
	return newResponseData(resp.StatusCode, r, resp.Header, err), nil
}

func callRequestBodyString[T ResponseBody](ctx context.Context, method string, client *RestClient, path string, bodyString string, headers map[string]string) ResponseData[T] {
	if !client.cb.Ready() {
		return newResponseData[T](http.StatusInternalServerError, nil, nil, ErrServiceNotAvailable)
	}

	var err error
	defer func() {
		err = client.cb.Done(ctx, err)
	}()

	txn := monitoring.GetTransactionInContext(ctx)
	if txn != nil {
		segment := createSegment(ctx, txn, method, path)
		defer monitoring.EndTransactionSegment(segment)
	}

	dataIOReader := strings.NewReader(bodyString)

	url := fmt.Sprintf("%s%s", client.baseURL, path)

	req, err := http.NewRequest(method, url, dataIOReader)
	if err != nil {
		return newResponseData[T](http.StatusInternalServerError, nil, nil, fmt.Errorf("could not create new request: %w", err))
	}

	addHeaders(req, headers)

	resp, err := client.client.Do(req)
	if err != nil {
		return newResponseData[T](http.StatusInternalServerError, nil, nil, err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			logging.Error("could not close response body: %v", err)
		}
	}()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return newResponseData[T](http.StatusInternalServerError, nil, nil, err)
	}

	var r T
	if err := json.Unmarshal(bodyText, &r); err != nil {
		return newResponseData[T](resp.StatusCode, nil, resp.Header, err)
	}
	return newResponseData(resp.StatusCode, &r, resp.Header, err)
}

func createSegment(ctx context.Context, txn interface{}, method string, path string) interface{} {
	segment := monitoring.StartTransactionSegment(ctx, txn, restClientTransaction, map[string]any{
		"method": method,
		"path":   path,
	})
	return segment
}

func processErrorResponse[E ResponseErrorData](resp *http.Response) (*E, error) {
	respErr, derr := decodeResponseErrorData[E](resp)
	if respErr != nil {
		return respErr, fmt.Errorf("%d statusCode", resp.StatusCode)
	} else if derr != nil {
		return nil, fmt.Errorf("%d statusCode. Body decoder Error: %w", resp.StatusCode, derr)
	}
	return nil, fmt.Errorf("%d statusCode", resp.StatusCode)
}

func makeBytesBody(entityBody any) (*bytes.Buffer, error) {
	if entityBody == nil {
		return nil, nil
	}
	requestBody, err := json.Marshal(entityBody)
	if err != nil {
		return nil, fmt.Errorf("could not marshal entity body: %w", err)
	}
	return bytes.NewBuffer(requestBody), nil
}

func decodeResponse[T ResponseBody](resp *http.Response) (*T, error) {
	defer func() {
		closeBody(resp)
	}()
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	var responseModel T
	err := json.NewDecoder(resp.Body).Decode(&responseModel)
	if err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}
	return &responseModel, nil
}

func decodeResponseErrorData[E ResponseErrorData](resp *http.Response) (*E, error) {
	defer func() {
		closeBody(resp)
	}()
	var responseModel E
	err := json.NewDecoder(resp.Body).Decode(&responseModel)
	switch {
	case err == io.EOF:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("could not decode response: %w", err)
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
