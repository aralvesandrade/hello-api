package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var err error

	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
		slog.Error("could not get hostname")
	}

	logLevel := parseLogLevel(os.Getenv("LOGGING_LEVEL"))

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: false,
	}).WithAttrs([]slog.Attr{
		slog.Group("app_details",
			slog.Int("pid", os.Getpid()),
			slog.String("hostname", hostname),
			slog.String("go_version", runtime.Version()),
		),
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// _, err = database.InitDbPostgres()
	// if err != nil {
	// 	slog.Error(err.Error())
	// }

	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/ping", pingHandler)
	http.Handle("/metrics", promhttp.Handler())

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 5001
	}

	slog.Info("server listening", "port", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		slog.Error(err.Error())
	}
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	msg := "Hello World!"
	slog.Debug(msg)
	w.Write([]byte(msg))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	msg := "up"
	slog.Debug(msg)
	w.Write([]byte(msg))
}

func parseLogLevel(logLevel string) slog.Level {
	logLevel = strings.ToUpper(logLevel)
	switch logLevel {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
