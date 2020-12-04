package internal

import (
	"fmt"
	"net/http"
	"os/user"
	"time"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// Log is the logging object used by the binary.
var Log *logrus.Logger = newLogger()

// HTTPLogger returns a http.Handler that logs the http request.
func HTTPLogger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		Log.Infof(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

// TODO(dotslash): Currently running tests also creates log files in ~/log
// directory. This needs to be fixed.
func newLogger() *logrus.Logger {
	fmt.Println("getting logger")
	start := time.Now()
	usr, _ := user.Current()
	home := usr.HomeDir
	lLog := logrus.New()
	lLog.Hooks.Add(lfshook.NewHook(
		lfshook.PathMap{
			logrus.InfoLevel:  home + "/log/bm-info.log",
			logrus.ErrorLevel: home + "/log/bm-error.log",
		},
		&logrus.TextFormatter{}))
	lLog.Info("got logger in ", time.Since(start))
	return lLog
}
