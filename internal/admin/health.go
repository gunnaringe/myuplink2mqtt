package admin

import (
	"fmt"
	"net/http"
)

var status = false

func Healthy() {
	status = true
}

func Unhealthy() {
	status = false
}

func healthy(w http.ResponseWriter, req *http.Request) {
	if status {
		fmt.Fprintf(w, "healthy")
	} else {
		w.WriteHeader(500)
		fmt.Fprintf(w, "unhealthy")
	}
}

func Start() {
	go func() {
		http.HandleFunc("/", healthy)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}()
}
