package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logLevel := parseLogLevel(os.Getenv("LOGGING_LEVEL"))
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/ping", pingHandler)
	http.Handle("/metrics", promhttp.Handler())

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 5001
	}

	slog.Info("server listening", "port", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
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
