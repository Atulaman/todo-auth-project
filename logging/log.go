package log

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// func Init() {
// 	Logger = logrus.New()
// 	Logger.SetFormatter(&logrus.JSONFormatter{
// 		FieldMap: logrus.FieldMap{
// 			logrus.FieldKeyTime:  "timestamp",
// 			logrus.FieldKeyLevel: "severity",
// 			logrus.FieldKeyMsg:   "message",
// 			logrus.FieldKeyFunc:  "func",
// 			logrus.FieldKeyFile:  "file",
// 		},
// 	})
// 	Logger.SetReportCaller(true)
// 	Logger.SetLevel(logrus.InfoLevel)
// 	//Logger.SetOutput(log.StandardLogger().Out)
// 	Logger.SetOutput(os.Stdout)
// }

func init() {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			//logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc: "func",
			logrus.FieldKeyFile: "file",
		},
	})
	Logger.SetOutput(os.Stdout)
}

func Logging(err error, message string, code int, severity string, r *http.Request) {
	var (
		path   string
		method string
	)
	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
	}
	if r != nil {
		path = r.URL.Path
		method = r.Method
	}
	entity := Logger.WithFields(logrus.Fields{
		//"message": message,
		"Path":       path,
		"Method":     method,
		"StatusCode": code,
		"Error":      errorMessage,
	})
	switch severity {
	case "info":
		entity.Info(message)
	case "warning":
		entity.Warn(message)
	case "error":
		entity.Error(message + " " + err.Error())
	default:
		entity.Fatal(message + " " + err.Error())
	}
}
