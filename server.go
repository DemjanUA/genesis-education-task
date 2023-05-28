package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	api_verions := "0.1"

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// /api
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Api verion: %s\n", api_verions)
	})

	api.HandleFunc("/rate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "success",
			"statusCode": 200,
			"data":       "/rate",
		})
	}).Methods("GET")

	api.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "success",
			"statusCode": 200,
			"data":       "/subscribe",
		})
	}).Methods("POST")

	api.HandleFunc("/sendEmails", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "success",
			"statusCode": 200,
			"data":       "/sendEmails",
		})
	}).Methods("POST")

	http.ListenAndServe(":80", r)
}
