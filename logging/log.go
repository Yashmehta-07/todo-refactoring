package logging

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// Initialize the Logger
var Logger = logrus.New()

func init() {
	// Configure logrus
	Logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyFunc:  "func",
			logrus.FieldKeyFile:  "file",
		},
	}) // JSON format for structured logging
	Logger.SetOutput(os.Stdout) // Output to standard output
}

func Log(err error, message string, severity string, code int, r *http.Request) {

	// Check if r is nil and provide default values if so
	var method, path, error string

	if err != nil {
		error = err.Error()
	}

	if r != nil {
		method = r.Method
		path = r.RequestURI
	}

	// Set the severity level dynamically
	entry := Logger.WithFields(logrus.Fields{
		"method":     method,
		"path":       path,
		"statusCode": code,
		"error":      error,
		// "message":  message,
	})

	// Log based on severity
	switch severity {
	case "info":
		entry.Info(message)
	case "warning":
		entry.Warn(message)
	case "error":
		entry.Error(message)
	case "fatal":
		entry.Fatal(message)
	case "debug":
		entry.Debug(message)
	default:
		entry.Info(message) // Default to info if severity is invalid
	}
}
