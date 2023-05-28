package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type CoinlayerResponse struct {
	Target string `json:"target"`
	Rates  struct {
		BTC float64 `json:"BTC"`
	} `json:"rates"`
}

func main() {
	api_verions := "0.1"

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Api verion: %s\n", api_verions)
	})

	api.HandleFunc("/rate", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://api.coinlayer.com/live?access_key=06f05f91a0c78ceb874adc4d6e65bdb2&target=UAH")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		// fmt.Println(string(body))

		var result CoinlayerResponse
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}
		fmt.Println(PrettyPrint(result))

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "success",
			"statusCode": 200,
			"value":      result.Rates.BTC,
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

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
