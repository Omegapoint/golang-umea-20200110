package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		c, _ := newClient("192.168.1.1", 8181, "Sample client", time.Now())
		json.NewEncoder(w).Encode(c)
	})
	_ = http.ListenAndServe(":8181", nil)
}
