/*
Copyright 2020 The go-harbor Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
*/

package rest

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	flowcontrol2 "github.com/TimeBye/go-harbor/pkg/rest/util/flowcontrol"
	"go/types"
	"golang.org/x/net/http2"
	"io"
	"io/ioutil"
	"k8s.io/klog"
	"mime"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	// longThrottleLatency defines threshold for logging requests. All requests being
	// throttle for more than longThrottleLatency will be logged.
	longThrottleLatency = 50 * time.Millisecond
)

// HTTPClient is an interface for testing a request object.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ResponseWrapper is an interface for getting a response.
// The response may be either accessed as a raw data (the whole output is put into memory) or as a stream.
type ResponseWrapper interface {
	DoRaw() ([]byte, error)
	Stream() (io.ReadCloser, error)
}

// RequestConstructionError is returned when there's an error assembling a request.
type RequestConstructionError struct {
	Err error
}

// Error returns a textual description of 'r'.
func (r *RequestConstructionError) Error() string {
	return fmt.Sprintf("request construction error: '%v'", r.Err)
}

// Request allows for building up a request to a server in a chained fashion.
// Any errors are stored until the end of your call, so you only have to
// check once.
type Request struct {
	// required
	client HTTPClient
	verb   string

	baseURL *url.URL
	//serializers Serializers
	// generic components accessible via method setters
	pathPrefix string
	subpath    string
	params     url.Values
	headers    http.Header

	resource     string
	resourceName string
	subresource  string
	timeout      time.Duration
	project      string
	projectSet   bool
	// output
	err     error
	body    io.Reader
	content ContentConfig
	// This is only used for per-request timeouts, deadlines, and cancellations.
	ctx context.Context

	throttle flowcontrol2.RateLimiter
}

// Result contains the result of calling Request.Do().
type Result struct {
	body        []byte
	contentType string
	err         error
	statusCode  int
}

type ContentConfig struct {
	// AcceptContentTypes specifies the types the client will accept and is optional.
	// If not set, ContentType will be used to define the Accept header
	AcceptContentTypes string
	// ContentType specifies the wire format used to communicate with the server.
	// This value will be set as the Accept header on requests made to the server, and
	// as the default content type on any object sent to the server. If not set,
	// "application/json" is used.
	ContentType string
}

// NewRequest creates a new request helper object for accessing runtime.Objects on a server.
func NewRequest(client HTTPClient, verb string, baseURL *url.URL, headers map[string]string, versionedAPIPath string, content ContentConfig, throttle flowcontrol2.RateLimiter, timeout time.Duration) *Request {
	if headers == nil {
		headers = make(map[string]string, 0)
	}
	pathPrefix := "/"
	if baseURL != nil {
		pathPrefix = path.Join(pathPrefix, baseURL.Path)
	}
	r := &Request{
		client:     client,
		verb:       verb,
		baseURL:    baseURL,
		pathPrefix: path.Join(pathPrefix, versionedAPIPath),
		content:    content,
		throttle:   throttle,
		timeout:    timeout,
	}
	for k, v := range headers {
		r.SetHeader(k, v)
	}
	switch {
	case len(content.AcceptContentTypes) > 0:
		r.SetHeader("Accept", content.AcceptContentTypes)
	case len(content.ContentType) > 0:
		r.SetHeader("Accept", content.ContentType+", */*")
	}
	return r
}

// Prefix adds segments to the relative beginning to the request path. These
// items will be placed before the optional Namespace, Resource, or Name sections.
// Setting AbsPath will clear any previously set Prefix segments
func (r *Request) Prefix(segments ...string) *Request {
	if r.err != nil {
		return r
	}
	r.pathPrefix = path.Join(r.pathPrefix, path.Join(segments...))
	return r
}

// Suffix appends segments to the end of the path. These items will be placed after the prefix and optional
// Namespace, Resource, or Name sections.
func (r *Request) Suffix(segments ...string) *Request {
	if r.err != nil {
		return r
	}
	r.subpath = path.Join(r.subpath, path.Join(segments...))
	return r
}

func (r *Request) setParam(paramName, value string) *Request {
	if r.params == nil {
		r.params = make(url.Values)
	}
	r.params[paramName] = append(r.params[paramName], value)
	return r
}

func (r *Request) SetHeader(key string, values ...string) *Request {
	if r.headers == nil {
		r.headers = http.Header{}
	}
	r.headers.Del(key)
	for _, value := range values {
		r.headers.Add(key, value)
	}
	return r
}

// Timeout makes the request use the given duration as an overall timeout for the
// request. Additionally, if set passes the value as "timeout" parameter in URL.
func (r *Request) Timeout(d time.Duration) *Request {
	if r.err != nil {
		return r
	}
	r.timeout = d
	return r
}

// Body makes the request use obj as the body. Optional.
// If obj is a string, try to read a file of that name.
// If obj is a []byte, send it directly.
// If obj is an io.Reader, use it directly.
// If obj is a runtime.Object, marshal it correctly, and set Content-Type header.
// If obj is a runtime.Object and nil, do nothing.
// Otherwise, set an error.
func (r *Request) Body(obj interface{}) *Request {
	if r.err != nil {
		return r
	}
	switch t := obj.(type) {
	case string:
		data, err := ioutil.ReadFile(t)
		if err != nil {
			r.err = err
			return r
		}
		glogBody("Request Body", data)
		r.body = bytes.NewReader(data)
	case []byte:
		glogBody("Request Body", t)
		r.body = bytes.NewReader(t)
	case io.Reader:
		r.body = t
	case types.Struct:
		if reflect.ValueOf(t).IsNil() {
			return r
		}
		data, _ := json.Marshal(t)
		r.body = bytes.NewReader(data)
		r.SetHeader("Content-Type", r.content.ContentType)
	case interface{}:
		data, _ := json.Marshal(t)
		r.body = bytes.NewReader(data)
		r.SetHeader("Content-Type", r.content.ContentType)
		/*	case runtime.Object:
			// callers may pass typed interface pointers, therefore we must check nil with reflection
			if reflect.ValueOf(t).IsNil() {
				return r
			}
			data, err := runtime.Encode(r.serializers.Encoder, t)
			if err != nil {
				r.err = err
				return r
			}
			glogBody("Request Body", data)
			r.body = bytes.NewReader(data)
			r.SetHeader("Content-Type", r.content.ContentType)*/
	default:
		r.err = fmt.Errorf("unknown type used for body: %+v", obj)
	}
	return r
}

// Param creates a query parameter with the given string value.
func (r *Request) Param(paramName, s string) *Request {
	if r.err != nil {
		return r
	}
	return r.setParam(paramName, s)
}

// RequestURI overwrites existing path and parameters with the value of the provided server relative
// URI.
func (r *Request) RequestURI(uri string) *Request {
	if r.err != nil {
		return r
	}
	locator, err := url.Parse(uri)
	if err != nil {
		r.err = err
		return r
	}
	r.pathPrefix = locator.Path
	if len(locator.Query()) > 0 {
		if r.params == nil {
			r.params = make(url.Values)
		}
		for k, v := range locator.Query() {
			r.params[k] = v
		}
	}
	return r
}

// glogBody logs a body output that could be either JSON or protobuf. It explicitly guards against
// allocating a new string for the body output unless necessary. Uses a simple heuristic to determine
// whether the body is printable.
func glogBody(prefix string, body []byte) {
	if klog.V(8) {
		if bytes.IndexFunc(body, func(r rune) bool {
			return r < 0x0a
		}) != -1 {
			klog.Infof("%s:\n%s", prefix, truncateBody(hex.Dump(body)))
		} else {
			klog.Infof("%s: %s", prefix, truncateBody(string(body)))
		}
	}
}

// AbsPath overwrites an existing path with the segments provided. Trailing slashes are preserved
// when a single segment is passed.
func (r *Request) AbsPath(segments ...string) *Request {
	if r.err != nil {
		return r
	}
	r.pathPrefix = path.Join(r.baseURL.Path, path.Join(segments...))
	if len(segments) == 1 && (len(r.baseURL.Path) > 1 || len(segments[0]) > 1) && strings.HasSuffix(segments[0], "/") {
		// preserve any trailing slashes for legacy behavior
		r.pathPrefix += "/"
	}
	return r
}

// URL returns the current working URL.
func (r *Request) URL() *url.URL {
	p := r.pathPrefix
	if r.projectSet && len(r.project) > 0 {
		p = path.Join(p, "projects", r.project)
	}
	if len(r.resource) != 0 {
		p = path.Join(p, strings.ToLower(r.resource))
	}
	// Join trims trailing slashes, so preserve r.pathPrefix's trailing slash for backwards compatibility if nothing was changed
	if len(r.resourceName) != 0 || len(r.subpath) != 0 || len(r.subresource) != 0 {
		p = path.Join(p, r.resourceName, r.subresource, r.subpath)
	}

	finalURL := &url.URL{}
	if r.baseURL != nil {
		*finalURL = *r.baseURL
	}
	finalURL.Path = p

	query := url.Values{}
	for key, values := range r.params {
		for _, value := range values {
			query.Add(key, value)
		}
	}

	// timeout is handled specially here.
	if r.timeout != 0 {
		query.Set("timeout", r.timeout.String())
	}
	finalURL.RawQuery = query.Encode()
	return finalURL
}

// Project applies the namespace scope to a request (<resource>/[ns/<namespace>/]<name>)
func (r *Request) Project(project string) *Request {
	if r.err != nil {
		return r
	}
	if r.projectSet {
		r.err = fmt.Errorf("namespace already set to %q, cannot change to %q", r.project, project)
		return r
	}
	r.projectSet = true
	r.project = project
	return r
}

// Resource sets the resource to access (<resource>/[ns/<namespace>/]<name>)
func (r *Request) Resource(resource string) *Request {
	if r.err != nil {
		return r
	}
	if len(r.resource) != 0 {
		r.err = fmt.Errorf("resource already set to %q, cannot change to %q", r.resource, resource)
		return r
	}
	r.resource = resource
	return r
}

// Name sets the name of a resource to access (<resource>/[ns/<namespace>/]<name>)
func (r *Request) Name(resourceName string) *Request {
	if r.err != nil {
		return r
	}
	if len(resourceName) == 0 {
		r.err = fmt.Errorf("resource name may not be empty")
		return r
	}
	if len(r.resourceName) != 0 {
		r.err = fmt.Errorf("resource name already set to %q, cannot change to %q", r.resourceName, resourceName)
		return r
	}
	r.resourceName = resourceName
	return r
}

func (r *Request) tryThrottle() error {
	if r.throttle == nil {
		return nil
	}

	now := time.Now()
	var err error
	r.throttle.Accept()
	if latency := time.Since(now); latency > longThrottleLatency {
		klog.V(4).Infof("Throttling request took %v, request: %s:%s", latency, r.verb, r.URL().String())
	}

	return err
}

// Do formats and executes the request. Returns a Result object for easy response
// processing.
//
// Error type:
//  * If the request can't be constructed, or an error happened earlier while building its
//    arguments: *RequestConstructionError
//  * If the server responds with a status: *errors.StatusError or *errors.UnexpectedObjectError
//  * http.Client.Do errors are returned directly.
func (r *Request) Do() Result {
	if err := r.tryThrottle(); err != nil {
		return Result{err: err}
	}

	var result Result
	err := r.request(func(req *http.Request, resp *http.Response) {
		result = r.transformResponse(resp, req)
	})
	if err != nil {
		return Result{err: err}
	}
	return result
}

// request connects to the server and invokes the provided function when a server response is
// received. It handles retry behavior and up front validation of requests. It will invoke
// fn at most once. It will return an error if a problem occurred prior to connecting to the
// server - the provided function is responsible for handling server errors.
func (r *Request) request(fn func(*http.Request, *http.Response)) error {
	if r.err != nil {
		klog.V(4).Infof("Error in request: %v", r.err)
		return r.err
	}

	// TODO: added to catch programmer errors (invoking operations with an object with an empty namespace)
	if (r.verb == "GET" || r.verb == "PUT" || r.verb == "DELETE") && r.projectSet && len(r.resourceName) > 0 && len(r.project) == 0 {
		return fmt.Errorf("an empty namespace may not be set when a resource name is provided")
	}
	if (r.verb == "POST") && r.projectSet && len(r.project) == 0 {
		return fmt.Errorf("an empty namespace may not be set during creation")
	}

	client := r.client
	if client == nil {
		client = http.DefaultClient
	}

	// Right now we make about ten retry attempts if we get a Retry-After response.
	maxRetries := 10
	retries := 0
	for {
		url := r.URL().String()
		req, err := http.NewRequest(r.verb, url, r.body)
		if err != nil {
			return err
		}
		if r.timeout > 0 {
			if r.ctx == nil {
				r.ctx = context.Background()
			}
			var cancelFn context.CancelFunc
			r.ctx, cancelFn = context.WithTimeout(r.ctx, r.timeout)
			defer cancelFn()
		}
		if r.ctx != nil {
			req = req.WithContext(r.ctx)
		}
		req.Header = r.headers

		if retries > 0 {
			// We are retrying the request that we already send to apiserver
			// at least once before.
			// This request should also be throttled with the client-internal throttler.
			if err := r.tryThrottle(); err != nil {
				return err
			}
		}
		resp, err := client.Do(req)
		if err != nil {
			// For the purpose of retry, we set the artificial "retry-after" response.
			// TODO: Should we clean the original response if it exists?
			resp = &http.Response{
				StatusCode: http.StatusInternalServerError,
				Header:     http.Header{"Retry-After": []string{"1"}},
				Body:       ioutil.NopCloser(bytes.NewReader([]byte{})),
			}
		}

		done := func() bool {
			// Ensure the response body is fully read and closed
			// before we reconnect, so that we reuse the same TCP
			// connection.
			defer func() {
				const maxBodySlurpSize = 2 << 10
				if resp.ContentLength <= maxBodySlurpSize {
					io.Copy(ioutil.Discard, &io.LimitedReader{R: resp.Body, N: maxBodySlurpSize})
				}
				resp.Body.Close()
			}()

			retries++
			if seconds, wait := checkWait(resp); wait && retries < maxRetries {
				if seeker, ok := r.body.(io.Seeker); ok && r.body != nil {
					_, err := seeker.Seek(0, 0)
					if err != nil {
						klog.V(4).Infof("Could not retry request, can't Seek() back to beginning of body for %T", r.body)
						fn(req, resp)
						return true
					}
				}

				klog.V(4).Infof("Got a Retry-After %ds response for attempt %d to %v", seconds, retries, url)
				return false
			}
			fn(req, resp)
			return true
		}()
		if done {
			return nil
		}
	}
}

// checkWait returns true along with a number of seconds if the server instructed us to wait
// before retrying.
func checkWait(resp *http.Response) (int, bool) {
	switch r := resp.StatusCode; {
	// any 500 error code and 429 can trigger a wait
	case r == http.StatusTooManyRequests, r >= 500:
	default:
		return 0, false
	}
	i, ok := retryAfterSeconds(resp)
	return i, ok
}

// retryAfterSeconds returns the value of the Retry-After header and true, or 0 and false if
// the header was missing or not a valid number.
func retryAfterSeconds(resp *http.Response) (int, bool) {
	if h := resp.Header.Get("Retry-After"); len(h) > 0 {
		if i, err := strconv.Atoi(h); err == nil {
			return i, true
		}
	}
	return 0, false
}

// DoRaw executes the request but does not process the response body.
func (r *Request) DoRaw() ([]byte, error) {
	if err := r.tryThrottle(); err != nil {
		return nil, err
	}

	var result Result
	err := r.request(func(req *http.Request, resp *http.Response) {
		result.body, result.err = ioutil.ReadAll(resp.Body)
		glogBody("Response Body", result.body)
		if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusPartialContent {
			result.err = r.transformUnstructuredResponseError(resp, req, result.body)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.body, result.err
}

// transformUnstructuredResponseError handles an error from the server that is not in a structured form.
// It is expected to transform any response that is not recognizable as a clear server sent error from the
// K8S API using the information provided with the request. In practice, HTTP proxies and client libraries
// introduce a level of uncertainty to the responses returned by servers that in common use result in
// unexpected responses. The rough structure is:
//
// 1. Assume the server sends you something sane - JSON + well defined error objects + proper codes
//    - this is the happy path
//    - when you get this output, trust what the server sends
// 2. Guard against empty fields / bodies in received JSON and attempt to cull sufficient info from them to
//    generate a reasonable facsimile of the original failure.
//    - Be sure to use a distinct error type or flag that allows a client to distinguish between this and error 1 above
// 3. Handle true disconnect failures / completely malformed data by moving up to a more generic client error
// 4. Distinguish between various connection failures like SSL certificates, timeouts, proxy errors, unexpected
//    initial contact, the presence of mismatched body contents from posted content types
//    - Give these a separate distinct error type and capture as much as possible of the original message
//
// TODO: introduce transformation of generic http.Client.Do() errors that separates 4.
func (r *Request) transformUnstructuredResponseError(resp *http.Response, req *http.Request, body []byte) error {
	if body == nil && resp.Body != nil {
		if data, err := ioutil.ReadAll(&io.LimitedReader{R: resp.Body}); err == nil {
			body = data
		}
	}
	//retryAfter, _ := retryAfterSeconds(resp)
	return r.newUnstructuredResponseError(body, resp.StatusCode, req)
}

// newUnstructuredResponseError instantiates the appropriate generic error for the provided input. It also logs the body.
func (r *Request) newUnstructuredResponseError(body []byte, statusCode int, req *http.Request) error {
	return fmt.Errorf("%s url:%s StatusCode: %d", req.Method, req.URL.Path, statusCode)
}

// transformResponse converts an API response into a structured API object
func (r *Request) transformResponse(resp *http.Response, req *http.Request) Result {
	var body []byte
	if resp.Body != nil {
		data, err := ioutil.ReadAll(resp.Body)
		switch err.(type) {
		case nil:
			body = data
		case http2.StreamError:
			// This is trying to catch the scenario that the server may close the connection when sending the
			// response body. This can be caused by server timeout due to a slow network connection.
			// TODO: Add test for this. Steps may be:
			// 1. client-go (or kubectl) sends a GET request.
			// 2. Apiserver sends back the headers and then part of the body
			// 3. Apiserver closes connection.
			// 4. client-go should catch this and return an error.
			klog.V(2).Infof("Stream error %#v when reading response body, may be caused by closed connection.", err)
			streamErr := fmt.Errorf("stream error when reading response body, may be caused by closed connection. Please retry. Original error: %v", err)
			return Result{
				err: streamErr,
			}
		default:
			klog.Errorf("Unexpected error when reading response body: %v", err)
			unexpectedErr := fmt.Errorf("unexpected error when reading response body. Please retry. Original error: %v", err)
			return Result{
				err: unexpectedErr,
			}
		}
	}

	glogBody("Response Body", body)

	// verify the content type is accurate
	contentType := resp.Header.Get("Content-Type")

	switch {
	case resp.StatusCode == http.StatusSwitchingProtocols:
		// no-op, we've been upgraded
	case resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusPartialContent:
		// calculate an unstructured error from the response which the Result object may use if the caller
		// did not return a structured error.
		//retryAfter, _ := retryAfterSeconds(resp)
		err := r.newUnstructuredResponseError(body, resp.StatusCode, req)
		return Result{
			body:        body,
			contentType: contentType,
			statusCode:  resp.StatusCode,
			err:         err,
		}
	}

	return Result{
		body:        body,
		contentType: contentType,
		statusCode:  resp.StatusCode,
	}
}

// isTextResponse returns true if the response appears to be a textual media type.
func isTextResponse(resp *http.Response) bool {
	contentType := resp.Header.Get("Content-Type")
	if len(contentType) == 0 {
		return true
	}
	media, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	return strings.HasPrefix(media, "text/")
}

// truncateBody decides if the body should be truncated, based on the glog Verbosity.
func truncateBody(body string) string {
	max := 0
	switch {
	case bool(klog.V(10)):
		return body
	case bool(klog.V(9)):
		max = 10240
	case bool(klog.V(8)):
		max = 1024
	}

	if len(body) <= max {
		return body
	}

	return body[:max] + fmt.Sprintf(" [truncated %d chars]", len(body)-max)
}

// Into stores the result into obj, if possible. If obj is nil it is ignored.
// If the returned object is of type Status and has .Status != StatusSuccess, the
// additional information in Status will be used to enrich the error.
func (r Result) Into(obj interface{}) error {
	if r.err != nil {
		// Check whether the result has a Status object in the body and prefer that.
		if r.contentType == "application/json; charset=utf-8" {
			return fmt.Errorf("%v message:%s", r.err, string(r.body))
		}
		return r.err
	}
	return json.Unmarshal(r.body, obj)
}

func (r Result) Error() error {
	if r.err != nil {
		// Check whether the result has a Status object in the body and prefer that.
		if r.contentType == "application/json; charset=utf-8" {
			return fmt.Errorf("%v message:%s", r.err, string(r.body))
		}
		return r.err
	}
	return nil
}

func (r *Request) Params(o interface{}) *Request {
	return r.Query(o)
}

func (r *Request) Query(content interface{}) *Request {
	v := reflect.ValueOf(content)
	switch v.Kind() {
	case reflect.Ptr:
		panic("pointers are not allowed")
	case reflect.String:
		r.queryString(v.String())
	case reflect.Struct:
		r.queryStruct(v.Interface())
	case reflect.Map:
		r.queryMap(v.Interface())
	default:
	}
	return r
}

func (r *Request) queryString(content string) {
	var val map[string]string
	if err := json.Unmarshal([]byte(content), &val); err == nil {
		for k, v := range val {
			r.setParam(k, v)
		}
	} else {
		if queryData, err := url.ParseQuery(content); err == nil {
			for k, queryValues := range queryData {
				for _, queryValue := range queryValues {
					r.setParam(k, string(queryValue))
				}
			}
		} else {
			panic(err)
		}
		// TODO: need to check correct format of 'field=val&field=val&...'
	}
}

func (r *Request) queryStruct(content interface{}) {
	if marshalContent, err := json.Marshal(content); err != nil {
		panic(err)
	} else {
		var val map[string]interface{}
		if err := json.Unmarshal(marshalContent, &val); err != nil {
			panic(err)
		} else {
			for k, v := range val {
				k = strings.ToLower(k)
				var queryVal string
				switch t := v.(type) {
				case string:
					queryVal = t
				case float64:
					queryVal = strconv.FormatFloat(t, 'f', -1, 64)
				case time.Time:
					queryVal = t.Format(time.RFC3339)
				default:
					j, err := json.Marshal(v)
					if err != nil {
						continue
					}
					queryVal = string(j)
				}
				r.setParam(k, queryVal)
			}
		}
	}
}

func (r *Request) queryMap(content interface{}) {
	r.queryStruct(content)
}
