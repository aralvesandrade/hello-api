package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/ping", pingHandler)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 5001
	}

	fmt.Printf("Server listening on port %d...", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().Format(time.RFC3339)
	msg := fmt.Sprintf("Ping! Current Time: %s", currentTime)
	w.Write([]byte(msg))
}
