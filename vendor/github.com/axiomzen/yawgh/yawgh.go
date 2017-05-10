package yawgh

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// for speed, we should be using json Encode/Decode
// but for debugging, we need to io.readall

// Marshaler allows you to marshal into v given a content Type
type Marshaler interface {
	Marshal(v interface{}, contentType string) ([]byte, error)
}

// Unmarshaler allows you to unmarshal the body into v given the contentType
type Unmarshaler interface {
	Unmarshal(data []byte, v interface{}, contentType string) error
}

// RequestInterceptor allows you to intercept the request
// you can return an error in which case the request will stop
type RequestInterceptor interface {
	InterceptRequest(r *http.Request, body []byte, err error) error
}

// ResponseInterceptor allows you to intercept the response
// you can return an error in which case the response will stop
type ResponseInterceptor interface {
	InterceptResponse(r *http.Response, body []byte, contentType string) error
}

// HTTPVerb are the verbs allowed
type HTTPVerb string

const (
	// Get GET
	Get HTTPVerb = "GET"
	// Put PUT
	Put HTTPVerb = "PUT"
	// Post POST
	Post HTTPVerb = "POST"
	// Update UPDATE
	Update HTTPVerb = "UPDATE"
	// Delete DELETE
	Delete HTTPVerb = "DELETE"
	// Options OPTIONS
	Options HTTPVerb = "OPTIONS"
	// Patch PATCH
	Patch HTTPVerb = "PATCH"
	// Head HEAD
	Head HTTPVerb = "HEAD"
)

// Request is the request configuration struct
type Request struct {
	headers             map[string]string
	cookies             map[string]*http.Cookie
	urlParams           url.Values
	verb                HTTPVerb
	transport           string
	domainHost          string
	port                uint
	urlPath             []string
	contentType         string
	requestBody         interface{}
	responseBody        interface{}
	errorResponseBody   interface{}
	marshaler           Marshaler
	unmarshaler         Unmarshaler
	requestInterceptor  RequestInterceptor
	responseInterceptor ResponseInterceptor
	numRetries          int
}

var jsonMime = "application/json"

// New creates a new request
func New() *Request {
	return &Request{
		headers:     make(map[string]string),
		cookies:     make(map[string]*http.Cookie),
		urlParams:   url.Values{},
		contentType: jsonMime,
		numRetries:  10,
		urlPath:     []string{},
	}
}

// withVerb sets the verb
func (c *Request) withVerb(verb HTTPVerb, u string) *Request {
	c.URLComponent(u)
	c.verb = verb
	return c
}

// Get configures a typical GET request
func (c *Request) Get(u string) *Request {
	return c.withVerb(Get, u)
}

// Put configures a typical PUT request
func (c *Request) Put(u string) *Request {
	return c.withVerb(Put, u)
}

// Post configures a typical POST request
func (c *Request) Post(u string) *Request {
	return c.withVerb(Post, u)
}

// Update configures a typical UPDATE request
func (c *Request) Update(u string) *Request {
	return c.withVerb(Update, u)
}

// Delete configures a typical DELETE request
func (c *Request) Delete(u string) *Request {
	return c.withVerb(Delete, u)
}

// Options configures a typical OPTIONS request
func (c *Request) Options(u string) *Request {
	return c.withVerb(Options, u)
}

// Patch configures a typical PATCH request
func (c *Request) Patch(u string) *Request {
	return c.withVerb(Patch, u)
}

// Head configures a typical HEAD request
func (c *Request) Head(u string) *Request {
	return c.withVerb(Head, u)
}

// Header adds a header
// the key and value are removed if the value is ""
func (c *Request) Header(key, value string) *Request {
	if value == "" {
		delete(c.headers, key)
	} else {
		c.headers[key] = value
	}

	return c
}

// SetCookie adds a cookie
// the key and value are removed if the value is "", a default of 1 hour expiry is used
func (c *Request) SetCookie(key, value string) *Request {

	if value == "" {
		delete(c.cookies, key)
	} else {
		cookie := http.Cookie{
			Name:    key,
			Value:   value,
			Expires: time.Now().Add(time.Hour),
		}
		c.cookies[key] = &cookie
	}

	return c
}

// Cookie adds a http cookie
// replaces any existing cookies of the same name, as it uses the cookie.Name as the key
func (c *Request) Cookie(cookie *http.Cookie) *Request {

	if cookie != nil {
		c.cookies[cookie.Name] = cookie
	}

	return c
}

// URLParam adds a url param
func (c *Request) URLParam(key, value string) *Request {
	c.urlParams.Add(key, value)
	return c
}

// URLComponent adds the url component to the url
func (c *Request) URLComponent(comp string) *Request {
	for _, s := range strings.Split(comp, "/") {
		if len(s) > 0 {
			c.urlPath = append(c.urlPath, s)
		}
	}
	return c
}

// RequestBody allows you to set request body of the request
func (c *Request) RequestBody(body interface{}) *Request {
	c.requestBody = body
	return c
}

// ResponseBody allows you to set the target struct to marshal the body of the
// response into
func (c *Request) ResponseBody(resp interface{}) *Request {
	c.responseBody = resp
	return c
}

// ErrorResponseBody allows you to set the target struct to marshal the body of the
// response into
func (c *Request) ErrorResponseBody(resp interface{}) *Request {
	c.errorResponseBody = resp
	return c
}

// Transport allows you to set the transport mechanism
func (c *Request) Transport(t string) *Request {
	c.transport = t
	return c
}

// DomainHost allows you to set the domain
func (c *Request) DomainHost(dh string) *Request {
	c.domainHost = dh
	return c
}

// Port allows you to set the domain
func (c *Request) Port(p uint) *Request {
	c.port = p
	return c
}

// ContentType allows you to set the content type and body of the request
// and an encoder/decoder
func (c *Request) ContentType(ct string) *Request {
	c.contentType = ct
	return c
}

// Marshaler allows you to set the marshaler for the request
func (c *Request) Marshaler(m Marshaler) *Request {
	c.marshaler = m
	return c
}

// Unmarshaler allows you to set the marshaler for the request
func (c *Request) Unmarshaler(u Unmarshaler) *Request {
	c.unmarshaler = u
	return c
}

// RequestInterceptor allows you to set the request interceptor for the request
func (c *Request) RequestInterceptor(r RequestInterceptor) *Request {
	c.requestInterceptor = r
	return c
}

// ResponseInterceptor allows you to set the marshaler for the request
func (c *Request) ResponseInterceptor(r ResponseInterceptor) *Request {
	c.responseInterceptor = r
	return c
}

// Retries allows you to set the marshaler for the request
func (c *Request) Retries(r int) *Request {
	c.numRetries = r
	return c
}

// checkAndClose private helper function
func checkAndClose(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}

// GetURL returns a full url based off of everything
func (c *Request) computeURL() string {
	var buffer bytes.Buffer
	buffer.WriteString(c.transport)
	buffer.WriteString("://")
	buffer.WriteString(c.domainHost)
	if string(c.domainHost[len(c.domainHost)-1]) != ":" {
		buffer.WriteString(":")
	}
	if c.port == 0 {
		// assume 80
		buffer.WriteString("80")
	} else {
		buffer.WriteString(strconv.FormatUint(uint64(c.port), 10))
	}
	buffer.WriteString("/")
	buffer.WriteString(strings.Join(c.urlPath, "/"))
	if len(c.urlParams) > 0 {
		buffer.WriteString("?")
		buffer.WriteString(c.urlParams.Encode())
	}
	return buffer.String()
}

// createRequest private helper function
func (c *Request) createRequest() (*http.Request, []byte, error) {
	var req *http.Request
	var err error
	var data []byte
	if c.requestBody == nil {
		req, err = http.NewRequest(string(c.verb), c.computeURL(), nil)
		if err != nil {
			return req, data, err
		}
	} else {
		if c.marshaler == nil {
			return req, data, errors.New("Please set a marshaler")
		}
		data, err = c.marshaler.Marshal(c.requestBody, c.contentType)
		if err != nil {
			return req, data, err
		}
		//gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
		req, err = http.NewRequest(string(c.verb), c.computeURL(), bytes.NewBuffer(data))
		if err != nil {
			return req, data, err
		}
	}
	// TODO: we can probably keep alive, this will be more efficient
	req.Close = true

	// Add headers to clientReq
	for key, value := range c.headers {
		req.Header.Add(key, value)
	}

	// Add cookies to clientReq
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	req.Header.Add("Content-Type", c.contentType)

	return req, data, err
}

// Do makes the request and returns the status code
func (c *Request) Do() (int, error) {
	resp, err := c.DoResponse()
	if resp == nil {
		return 0, err
	}

	return resp.StatusCode, err
}

// DoResponse makes the request and returns its response
// TODO: refactor
func (c *Request) DoResponse() (*http.Response, error) {
	client := http.Client{}
	req, data, err := c.createRequest()
	// do interceptor here?
	if c.requestInterceptor != nil {
		err = c.requestInterceptor.InterceptRequest(req, data, err)
	}
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	defer checkAndClose(resp)
	repeatCount := 0
	for err != nil && strings.HasSuffix(err.Error(), " EOF") && repeatCount < c.numRetries {
		repeatCount = repeatCount + 1
		time.Sleep(500 * time.Millisecond)
		//log.WithError(err).Error("requestPage")
		req, data, err = c.createRequest()
		if c.requestInterceptor != nil {
			err = c.requestInterceptor.InterceptRequest(req, data, err)
		}
		if err != nil {
			return nil, err
		}
		resp, err = client.Do(req)
	}

	if err != nil || resp == nil {
		return nil, err
	}

	// check to see if error condition
	if c.responseBody != nil && resp.StatusCode != http.StatusNoContent && resp.StatusCode/100 == 2 {
		err = unmarshalHelper(c.unmarshaler, resp, c.responseBody, c.contentType, c.responseInterceptor)
	} else if resp.StatusCode > 399 && c.errorResponseBody != nil {
		err = unmarshalHelper(c.unmarshaler, resp, c.errorResponseBody, c.contentType, c.responseInterceptor)
	} else {
		data, err = ioutil.ReadAll(resp.Body)
		// do response interceptor here
		if c.responseInterceptor != nil {
			err = c.responseInterceptor.InterceptResponse(resp, data, c.contentType)
		}
	}
	return resp, err
}

// helper function
func unmarshalHelper(unmarshaler Unmarshaler, resp *http.Response, v interface{}, contentType string, respInterceptor ResponseInterceptor) error {
	if unmarshaler == nil {
		return errors.New("Please set an unmarshaler")
	}

	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// do response interceptor here
	if respInterceptor != nil {
		err = respInterceptor.InterceptResponse(resp, bb, contentType)
	}

	if err != nil {
		return err
	}

	return unmarshaler.Unmarshal(bb, v, contentType)
}
