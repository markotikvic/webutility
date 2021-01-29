package logger

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

const splitLine = "=============================================================="

// LogHTTPRequest ...
func (l *Logger) LogHTTPRequest(req *http.Request, userID string) string {
	if userID == "" {
		userID = "-"
	}

	var b strings.Builder

	b.WriteString("Request:\n")
	// CLF-like header
	fmt.Fprintf(&b, "%s %s %s\n", req.RemoteAddr, userID, time.Now().Format(dateTimeFormat))

	body, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	const sepStr = "\r\n\r\n"
	sepIndex := bytes.Index(body, []byte(sepStr))
	if sepIndex == -1 {
		b.WriteString(string(body) + "\n\n")
	} else {
		sepIndex += len(sepStr)
		payload, _ := printJSON(body[sepIndex:])
		b.WriteString(string(body[:sepIndex]) + string(payload) + "\n\n")
	}

	return b.String()
}

// LogHTTPResponse ...
func (l *Logger) LogHTTPResponse(data []byte, status int, startTime time.Time) string {
	duration := time.Now().Sub(startTime)
	jsonData, _ := printJSON(data)
	return fmt.Sprintf("Response:\n%d %v %dB\n%s\n%s\n\n", status, duration, len(data), jsonData, splitLine)
}

// CombineHTTPLogs ...
func (l *Logger) CombineHTTPLogs(in string, out string) {
	if l.outputFile == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	msg := in + out
	if l.shouldSplit(len(msg)) {
		l.split()
	}
	l.outputFile.WriteString(msg)
}
