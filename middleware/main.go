package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	web "git.to-net.rs/marko.tikvic/webutility"
)

func Headers(h http.HandlerFunc) http.HandlerFunc {
	return SetAccessControlHeaders(IgnoreOptionsRequests(ParseForm(h)))
}

func AuthUser(roles string, h http.HandlerFunc) http.HandlerFunc {
	return SetAccessControlHeaders(IgnoreOptionsRequests(ParseForm(Auth(roles, h))))
}

func AuthUserAndLog(roles string, h http.HandlerFunc) http.HandlerFunc {
	return SetAccessControlHeaders(IgnoreOptionsRequests(ParseForm(LogHTTP(Auth(roles, h)))))
}

func LogTraffic(h http.HandlerFunc) http.HandlerFunc {
	return SetAccessControlHeaders(IgnoreOptionsRequests(ParseForm(LogHTTP(h))))
}

func TrafficLogsHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		logfile := req.FormValue("logfile")
		if logfile == "" {
			files, err := ioutil.ReadDir(httpLogger.GetOutDir() + "/")
			if err != nil {
				web.InternalServerError(w, req, err.Error())
				return
			}

			errorLogs := make([]os.FileInfo, 0)
			httpLogs := make([]os.FileInfo, 0)

			var errorLogsCount, httpLogsCount int
			for _, f := range files {
				if strings.HasPrefix(f.Name(), "err") {
					errorLogs = append(errorLogs, f)
					errorLogsCount++
				} else if strings.HasPrefix(f.Name(), "http") {
					httpLogs = append(httpLogs, f)
					httpLogsCount++
				}
			}

			web.SetContentType(w, "text/html; charset=utf-8")
			web.SetResponseStatus(w, http.StatusOK)

			web.WriteResponse(w, []byte(`
				<body style='background-color: black; color: white'>
				<table>
				<tr>
				  <th>Error logs</th><th></th> <th style="width: 25px"></th>
				  <th>Traffic logs</th><th></th>
				</tr>
			`))

			var (
				div, name string
				size      int64
			)

			max := errorLogsCount
			if httpLogsCount > errorLogsCount {
				max = httpLogsCount
			}

			for i := 0; i < max; i++ {
				div = "<tr>"

				if i < errorLogsCount {
					name = errorLogs[i].Name()
					size = errorLogs[i].Size()
					div += fmt.Sprintf(`
						<td>
						  <a style="color: white"
						     href="/api/v1/logs?logfile=%s"
						     target="_blank">%s
						   </a>
					        </td>
						<td style="color: white; text-align:right">%dB</td>`,
						name, name, size,
					)
				} else {
					div += fmt.Sprintf(`<td></td><td></td>`)
				}

				div += "<td></td>"

				if i < httpLogsCount {
					name := httpLogs[i].Name()
					size := httpLogs[i].Size()
					div += fmt.Sprintf(`
						<td>
						  <a style="color: white"
						     href="/api/v1/logs?logfile=%s"
						     target="_blank">%s
						  </a></td>
						<td style="color: white; text-align:right">%dB</td>`,
						name, name, size,
					)
				} else {
					div += fmt.Sprintf(`<td></td><td></td>`)
				}

				div += "</tr>"
				web.WriteResponse(w, []byte(div))
			}
			web.WriteResponse(w, []byte("</table></body>"))
		} else {
			content, err := web.ReadFileContent(httpLogger.GetOutDir() + "/" + logfile)
			if err != nil {
				web.InternalServerError(w, req, err.Error())
				return
			}
			web.SetResponseStatus(w, http.StatusOK)
			web.WriteResponse(w, []byte("<body style='background-color: black; color: white'>"))
			web.WriteResponse(w, []byte("<pre>"))
			web.WriteResponse(w, content)
			web.WriteResponse(w, []byte("</pre></body>"))
		}
	}
}
