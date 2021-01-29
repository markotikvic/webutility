package middleware

import (
	"net/http"
	"time"

	web "git.to-net.rs/marko.tikvic/webutility"
	"git.to-net.rs/marko.tikvic/webutility/logger"
)

var httpLogger *logger.Logger

func SetAccessControlHeaders(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		web.SetAccessControlHeaders(w)

		h(w, req)
	}
}

// IgnoreOptionsRequests ...
func IgnoreOptionsRequests(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodOptions {
			return
		}

		h(w, req)
	}
}

// ParseForm ...
func ParseForm(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			web.BadRequest(w, req, err.Error())
			return
		}

		h(w, req)
	}
}

// ParseMultipartForm ...
func ParseMultipartForm(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseMultipartForm(32 << 20)
		if err != nil {
			web.BadRequest(w, req, err.Error())
			return
		}

		h(w, req)
	}
}

// SetLogger ...
func SetLogger(logger *logger.Logger) {
	httpLogger = logger
}

func StartLogging(filename, dir string) (err error) {
	if httpLogger, err = logger.New(filename, dir, logger.MaxLogSize1MB); err != nil {
		return err
	}
	return nil
}

func CloseLogger() {
	httpLogger.Close()
}

// LogHTTP ...
func LogHTTP(hfunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if httpLogger == nil {
			hfunc(w, req)
			return
		}

		t1 := time.Now()

		claims, _ := web.GetTokenClaims(req)
		in := httpLogger.LogHTTPRequest(req, claims.Username)

		rec := web.NewStatusRecorder(w)

		hfunc(rec, req)

		out := httpLogger.LogHTTPResponse(rec.Data(), rec.Status(), t1)

		httpLogger.CombineHTTPLogs(in, out)
	}
}

// Auth ...
func Auth(roles string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if _, err := web.AuthCheck(req, roles); err != nil {
			web.Unauthorized(w, req, err.Error())
			return
		}

		h(w, req)
	}
}
