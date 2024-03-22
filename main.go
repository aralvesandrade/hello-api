package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/ping", pingHandler)
	http.ListenAndServe(":5001", nil)
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().Format(time.RFC3339)
	msg := fmt.Sprintf("Ping! Current Time: %s", currentTime)
	w.Write([]byte(msg))
}
