package webutility

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// StatusRecorder ...
type StatusRecorder struct {
	writer http.ResponseWriter
	status int
	size   int
	data   []byte
}

// NewStatusRecorder ...
func NewStatusRecorder(w http.ResponseWriter) *StatusRecorder {
	return &StatusRecorder{
		writer: w,
		status: 0,
		size:   0,
		data:   nil,
	}
}

// WriteHeader is a wrapper http.ResponseWriter interface
func (r *StatusRecorder) WriteHeader(code int) {
	r.status = code
	r.writer.WriteHeader(code)
}

// Write is a wrapper for http.ResponseWriter interface
func (r *StatusRecorder) Write(in []byte) (int, error) {
	r.size = len(in)
	if r.status >= 400 {
		r.data = make([]byte, len(in))
		copy(r.data, in)
	}
	return r.writer.Write(in)
}

// Header is a wrapper for http.ResponseWriter interface
func (r *StatusRecorder) Header() http.Header {
	return r.writer.Header()
}

// Status ...
func (r *StatusRecorder) Status() int {
	return r.status
}

// Size ...
func (r *StatusRecorder) Size() int {
	return r.size
}

// Size ...
func (r *StatusRecorder) Data() []byte {
	return r.data
}

// NotFoundHandlerFunc writes HTTP error 404 to w.
func NotFoundHandlerFunc(w http.ResponseWriter, req *http.Request) {
	SetAccessControlHeaders(w)
	SetContentType(w, "application/json")
	NotFound(w, req, fmt.Sprintf("Resource you requested was not found: %s", req.URL.String()))
}

// SetContentType must be called before SetResponseStatus (w.WriteHeader) (?)
func SetContentType(w http.ResponseWriter, ctype string) {
	w.Header().Set("Content-Type", ctype)
}

// SetResponseStatus ...
func SetResponseStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

// WriteResponse ...
func WriteResponse(w http.ResponseWriter, content []byte) {
	w.Write(content)
}

// SetAccessControlHeaders set's default headers for an HTTP response.
func SetAccessControlHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// GetLocale ...
func GetLocale(req *http.Request, dflt string) string {
	loc := req.FormValue("locale")
	if loc == "" {
		return dflt
	}
	return loc
}

// Success ...
func Success(w http.ResponseWriter, payload interface{}, code int) {
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

// OK ...
func OK(w http.ResponseWriter, payload interface{}) {
	SetContentType(w, "application/json")
	Success(w, payload, http.StatusOK)
}

// Created ...
func Created(w http.ResponseWriter, payload interface{}) {
	SetContentType(w, "application/json")
	Success(w, payload, http.StatusCreated)
}

type weberror struct {
	Request string `json:"request"`
	Error   string `json:"error"`
	//Code    int64  `json:"code"` TODO
}

// Error ...
func Error(w http.ResponseWriter, r *http.Request, code int, err string) {
	werr := weberror{Error: err, Request: r.Method + " " + r.RequestURI}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(werr)
}

// BadRequest ...
func BadRequest(w http.ResponseWriter, r *http.Request, err string) {
	SetContentType(w, "application/json")
	Error(w, r, http.StatusBadRequest, err)
}

// Unauthorized ...
func Unauthorized(w http.ResponseWriter, r *http.Request, err string) {
	SetContentType(w, "application/json")
	Error(w, r, http.StatusUnauthorized, err)
}

// Forbidden ...
func Forbidden(w http.ResponseWriter, r *http.Request, err string) {
	SetContentType(w, "application/json")
	Error(w, r, http.StatusForbidden, err)
}

// NotFound ...
func NotFound(w http.ResponseWriter, r *http.Request, err string) {
	SetContentType(w, "application/json")
	Error(w, r, http.StatusNotFound, err)
}

// Conflict ...
func Conflict(w http.ResponseWriter, r *http.Request, err string) {
	SetContentType(w, "application/json")
	Error(w, r, http.StatusConflict, err)
}

// InternalServerError ...
func InternalServerError(w http.ResponseWriter, r *http.Request, err string) {
	SetContentType(w, "application/json")
	Error(w, r, http.StatusInternalServerError, err)
}

func SetHeader(r *http.Request, key, val string) {
	r.Header.Set(key, val)
}

func AddHeader(r *http.Request, key, val string) {
	r.Header.Add(key, val)
}

func GetHeader(r *http.Request, key string) string {
	return r.Header.Get(key)
}

func ClientUTCOffset(req *http.Request) int64 {
	return StringToInt64(GetHeader(req, "X-Timezone-Offset"))
}
