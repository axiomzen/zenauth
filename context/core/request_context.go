package core

import (
	"compress/flate"
	"compress/gzip"

	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/axiomzen/gorelic"
	"github.com/axiomzen/zenauth/config"
	"github.com/axiomzen/zenauth/constants"
	"github.com/axiomzen/zenauth/data"
	"github.com/axiomzen/zenauth/helpers"
	"github.com/axiomzen/zenauth/helpers/header"
	"github.com/axiomzen/zenauth/models"
	"github.com/gocraft/web"
	"github.com/newrelic/go-agent"
	"github.com/rcrowley/go-metrics"
	"github.com/twinj/uuid"
)

type (
	// tokenStatus is our token status states
	tokenStatus int

	// RequestContext will contain all of the data that will
	// flow through a request. This won't be persisted.
	RequestContext struct {
		requestID      uuid.UUID
		Log            *log.Entry
		statusCode     constants.HTTPStatusCode
		responseObject interface{}
		Config         *config.ZENAUTHConfig
		DAL            data.ZENAUTHProvider
		NewRelic       newrelic.Transaction
	}

	compressionResponseWriter struct {
		io.Writer
		web.ResponseWriter
		sniffDone bool
	}
)

const (
	tokenValid tokenStatus = iota
	tokenExpired
	tokenInvalid
	accessControlAllowOrigin      string = "Access-Control-Allow-Origin"
	accessControlAllowMethods            = "Access-Control-Allow-Methods"
	accessControlAllowHeaders            = "Access-Control-Allow-Headers"
	accessControlAllowCredentials        = "Access-Control-Allow-Credentials"
)

var (
	newRelicApp    *newrelic.Application
	newRelicPlugin *gorelic.Agent
)

func InitNewRelicApp(c *config.ZENAUTHConfig) bool {
	if c.NewRelicEnabled {
		cfg := newrelic.NewConfig(c.NewRelicName, c.NewRelicKey)
		app, err := newrelic.NewApplication(cfg)
		if err != nil {
			log.WithError(err).Errorf("Could not start New Relic monitoring")
			return false
		} else {
			newRelicApp = &app
			return true
		}
	}
	return false
}

// GoRelicHandler private wrapper for gorelic
func GoRelicHandler(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	defer newRelicPlugin.HTTPTimer.UpdateSince(time.Now())
	next(w, r)
}

// InitNewRelicPlugin private helper function for initialization of new relic handler
func InitNewRelicPlugin(c *config.ZENAUTHConfig) bool {

	if c.NewRelicEnabled {
		newRelicPlugin = gorelic.NewAgent()
		newRelicPlugin.Verbose = false
		newRelicPlugin.NewrelicName = c.NewRelicName
		newRelicPlugin.NewrelicLicense = c.NewRelicKey
		newRelicPlugin.CollectGcStat = c.NewRelicGCEnabled
		newRelicPlugin.HTTPTimer = metrics.NewTimer()
		newRelicPlugin.CollectHTTPStat = true
		newRelicPlugin.CollectMemoryStat = c.NewRelicMemEnabled
		newRelicPlugin.GCPollInterval = int(c.NewRelicGCPoll)
		newRelicPlugin.MemoryAllocatorPollInterval = int(c.NewRelicMemPoll)
		newRelicPlugin.NewrelicPollInterval = int(c.NewRelicPoll)
		if err := newRelicPlugin.Run(); err != nil {
			log.WithError(err).Errorf("Could not start New Relic plugin")
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}

// Setup creates a uuid for reach request, sets up the debug logger
func (c *RequestContext) Setup(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	// New UUID for logging & debugging purposes
	conf, _ := config.Get()
	c.Config = conf
	requestIDStr := r.Header.Get(conf.RequestIDHeader)
	if requestID, err := uuid.Parse(requestIDStr); err != nil {
		c.requestID = uuid.NewV4()
	} else {
		c.requestID = requestID
	}
	dal, _ := data.Get(conf)
	c.DAL = dal
	next(w, r)
}

// NewRelicTransaction starts and attaches a new relic agent transaction
func (c *RequestContext) NewRelicTransaction(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.NewRelic = (*newRelicApp).StartTransaction(r.Method+" "+r.RoutePath(), w, r.Request)
	defer c.NewRelic.End()

	next(w, r)

	// If there's a 5xx type error, log it explicitly
	switch c.statusCode {
	case constants.StatusInternalServerError, constants.StatusServiceUnavailable:

		if err, ok := c.responseObject.(error); ok {
			c.NewRelic.NoticeError(err)
		} else {
			c.NewRelic.NoticeError(errors.New("Internal Server Error"))
		}
	}
}

// Logging Middleware: logs incoming requests and responses
func (c *RequestContext) Logging(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	// setup the logger
	c.Log = log.WithField("server", "http").WithField("id", c.requestID.String())
	// add hostname
	c.Log = c.Log.WithField("host", r.Host)
	// Attach this id to the context with the corrosponding ID
	c.Log.WithFields(log.Fields{"method": r.Method, "path": r.URL.Path}).Debug("Received request")
	// call the next middleware
	next(w, r)

	// log the result
	switch c.statusCode {
	case constants.StatusFound, constants.StatusMovedPermanently:
		c.Log.WithFields(log.Fields{"location": w.Header().Get("Location"), "code": strconv.Itoa(int(c.statusCode))}).Debug("Redirected")
		break
	case constants.StatusInternalServerError,
		constants.StatusServiceUnavailable,
		constants.StatusUnauthorized,
		constants.StatusForbidden,
		constants.StatusExpiredToken,
		constants.StatusBadRequest:
		// see if this was an error, we could just return it here
		if eo, ok := (c.responseObject).(error); ok {
			c.Log.WithError(eo).WithField("code", strconv.Itoa(int(c.statusCode))).Error("Returned")
		} else {
			c.Log.WithError(errors.New("unknown error")).WithField("code", strconv.Itoa(int(c.statusCode))).Error("Returned")
		}
		break
	case constants.StatusNotFound:
		c.Log.WithFields(log.Fields{"method": r.Method, "code": strconv.Itoa(int(c.statusCode)), "path": r.URL.String()}).Warn("Not found")
		break
	case constants.StatusNoContent, constants.StatusCreated, constants.StatusOK:
		c.Log.WithFields(log.Fields{"method": r.Method, "code": strconv.Itoa(int(c.statusCode))}).Debug("Returned")
		break
	}
}

//AccessControlAllowHandler Middleware: handles preflight OPTIONS requests and sets up CORS headers
func (c *RequestContext) AccessControlAllowHandler(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	// w.Header().Set(AccessControlAllowOrigin, "*")
	w.Header().Set(accessControlAllowMethods, "POST, GET, OPTIONS, PUT, PATCH, DELETE")
	w.Header().Set(accessControlAllowHeaders, "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Origin, Accept, Authorization, "+c.Config.APITokenHeader+", "+c.Config.AuthTokenHeader)
	w.Header().Set(accessControlAllowCredentials, "true")

	next(w, r)
}

// OPTIONSHandler Middleware: return from OPTIONS requests as quickly as possible
func (c *RequestContext) OPTIONSHandler(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if r.Method == `OPTIONS` {
		w.WriteHeader(http.StatusOK)
		return
	}

	next(w, r)
}

// Write implemnents a writer
func (w *compressionResponseWriter) Write(b []byte) (int, error) {
	if !w.sniffDone {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", http.DetectContentType(b))
		}
		w.sniffDone = true
	}
	return w.Writer.Write(b)
}

// CompressionHandler Middleware: compresses the output if applicable
func (c *RequestContext) CompressionHandler(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	// see if they will accept gzip encoding
	//typical: "Accept-Encoding: gzip,deflate"

	w.Header().Add("Vary", "Accept-Encoding")

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		// only want to close if we don't panics
		defer func(gzw *gzip.Writer) {
			if r := recover(); r == nil {
				// no error, so flush buffer
				gzw.Close()
			} else {
				// throw that back up
				panic(r)
			}
		}(gz)
		// call next middleware
		next(&compressionResponseWriter{Writer: gz, ResponseWriter: w}, r)
	} else if strings.Contains(r.Header.Get("Accept-Encoding"), "deflate") {
		w.Header().Set("Content-Encoding", "deflate")
		f, _ := flate.NewWriter(w, int(c.Config.DeflateCompression))

		// only want to close if we don't panics
		defer func(fl *flate.Writer) {
			if r := recover(); r == nil {
				// no error, so flush buffer
				fl.Close()
			} else {
				// throw that back up
				panic(r)
			}
		}(f)
		next(&compressionResponseWriter{Writer: f, ResponseWriter: w}, r)
	} else {
		// compression not accepted, typically testing or debugging, etc.
		// TODO: perhaps a warning if we are not in the test environment
		next(w, r)
	}
}

// UnauthorizedHandler generic unauthorized handler
func (c *RequestContext) UnauthorizedHandler(w web.ResponseWriter, r *web.Request) {
	var model = models.NewErrorResponse(constants.APIUnauthorized, models.NewAZError("not authorized"))
	c.Render(constants.StatusUnauthorized, model, w, r)
}

// ExpiredHandler generic auth token expiry handler
func (c *RequestContext) ExpiredHandler(w web.ResponseWriter, r *web.Request) {
	var model = models.NewErrorResponse(constants.APIExpiredAuthToken, models.NewAZError("Expired auth token"))
	c.Render(constants.StatusExpiredToken, model, w, r)
}

// PingResponse Ping our webservice
//
// Type: GET
// Route: /ping
//
// Output:
//
//     HTTP 200
//       {
//         "ping": "pong"
//       }
func (c *RequestContext) PingResponse(rw web.ResponseWriter, req *web.Request) {

	if err := c.DAL.Ping(); err != nil {
		model := models.NewErrorResponse(constants.APIDatabaseUnreachable, models.NewAZError(err.Error()))
		c.Render(constants.StatusServiceUnavailable, model, rw, req)
		return
	}
	var ping models.Ping
	ping.Ping = "pong"
	c.Render(constants.StatusOK, &ping, rw, req)
}

// Render renders the interface to the response
func (c *RequestContext) Render(statusCode constants.HTTPStatusCode, v interface{}, w web.ResponseWriter, r *web.Request) {
	// all headers need to be set now before the call to WriteHeader
	// setting them after the call to WriteHeader doesn't make a difference
	// set the response object so we can log the error
	c.responseObject = v
	// take note of the response code so we know what we sent (for logging)
	c.statusCode = statusCode

	// this is an unresolved issue
	// the header needs to be written before the data is written
	// but if we have an error encoding something,
	// we have already written the header and cannot change it
	// we could use marshal instead of encode to render to memory first
	// however then we lose all advantages of streaming the response
	// so, hopefully we will catch any encoding issues via the tests

	// REMEMBER: for golang http, things must be in this order:
	// 1. w.Header().Set("key", "value")
	// 2. w.WriteHeader(200)
	// 3. w.Write(bytes)

	if v != nil {
		contentType := r.Header.Get("Content-Type")
		// check for presence of content-type
		if contentType == "" {
			// lets check accept header
			specs := header.ParseAccept(r.Header, "Accept")
			if len(specs) > 0 && len(specs[0].Value) > 0 {
				contentType = specs[0].Value
			} else {
				contentType = c.Config.DefaultContentType
			}
		}
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(int(c.statusCode))
		if err := helpers.Encode(v, contentType, w); err != nil {
			c.Log.Error(err)
		}
	} else {
		w.WriteHeader(int(c.statusCode))
	}
}

// DecodeHelper returns true upon sucessfully decoding the request into the interface
// otherwise it writes an error object to the response and returns false
func (c *RequestContext) DecodeHelper(v interface{}, message string, w web.ResponseWriter, r *web.Request) bool {

	if err := helpers.Decode(v, r.Header.Get("Content-Type"), r.Body); err != nil {
		model := models.NewErrorResponse(constants.APIParsing, models.NewAZError(err.Error()))
		// If we can't decode it, then it is probably a bad request (malformed json for example)
		c.Render(constants.StatusBadRequest, model, w, r)
		return false
	}
	return true
}

// Error Middleware: panics, etc
func (c *RequestContext) Error(rw web.ResponseWriter, r *web.Request, err interface{}) {
	// log here as the logging middleware is undone
	c.Log.Error(err)

	var errorObj *models.ErrorResponse

	if errorValue, ok := err.(error); ok {
		errorObj = models.NewErrorResponse(constants.APIPanic, models.NewAZError(errorValue.Error()))
	} else {
		errorObj = models.NewErrorResponse(constants.APIPanic, nil)
	}

	// at this point, we no longer have the compression middleware
	// so we have to redo it
	var nrw web.ResponseWriter
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		rw.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(rw)
		defer gz.Close()
		nrw = &compressionResponseWriter{Writer: gz, ResponseWriter: rw}
	} else if strings.Contains(r.Header.Get("Accept-Encoding"), "deflate") {
		rw.Header().Set("Content-Encoding", "deflate")
		flate, _ := flate.NewWriter(rw, int(c.Config.DeflateCompression))
		defer flate.Close()
		nrw = &compressionResponseWriter{Writer: flate, ResponseWriter: rw}
	} else {
		// compression not accepted, typically testing or debugging, etc.
		// TODO: perhaps a warning if we are not in the test environment
		nrw = rw
	}

	c.Render(constants.StatusInternalServerError, errorObj, nrw, r)
}

// NotFound Not Found function middleware
func (c *RequestContext) NotFound(rw web.ResponseWriter, req *web.Request) {
	if newRelicApp != nil {
		(*newRelicApp).RecordCustomEvent("Not Found", map[string]interface{}{"path": req.RequestURI})
	}
	model := models.NewErrorResponse(constants.APINotFound, nil)
	c.Render(constants.StatusNotFound, model, rw, req)
}
