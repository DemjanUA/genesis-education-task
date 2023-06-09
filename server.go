package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type CoinlayerResponse struct {
	Target string `json:"target"`
	Rates  struct {
		BTC float64 `json:"BTC"`
	} `json:"rates"`
}

type EmailStorage []string

func main() {
	api_verions := "0.1"

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Api verion: %s\n", api_verions)
	})

	api.HandleFunc("/rate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		BTCRate := getBTCValue(w)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "success",
			"statusCode": 200,
			"value":      BTCRate,
		})
	}).Methods("GET")

	api.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// ====== Dealing with www-form-urlencoded body ====== //
		r.ParseForm()

		userEmail, ok := r.PostForm["email"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":     "Bad Request",
				"statusCode": 400,
				"data":       "Invalid email parameter",
			})
			return
		}
		email := userEmail[0]

		// ====== Read file storage ====== //
		emailStorageFile, err := os.Open("email.storage.json")

		if err != nil {
			fmt.Println(err)
		}
		defer emailStorageFile.Close()

		byteEmailStorageFile, _ := ioutil.ReadAll(emailStorageFile)

		var emails EmailStorage
		if err := json.Unmarshal(byteEmailStorageFile, &emails); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}

		// ====== Check if email exist ====== //
		if arrayContains(emails, email) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":     "Conflict",
				"statusCode": 409,
				"data":       "Email address is already subscribed",
			})
			return
		}

		slice := append(emails, email)
		fmt.Println(PrettyPrint(slice))

		// ====== Write file storage ====== //
		content, err := json.Marshal(slice)
		err = ioutil.WriteFile("email.storage.json", content, 0644)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":     "Server Error",
				"statusCode": 500,
				"data":       "Internal Server Error",
			})
			return
		}

		// ====== Sending response ====== //
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "success",
			"statusCode": 200,
			"data":       "Subscribed",
		})
		return
	}).Methods("POST")

	api.HandleFunc("/sendEmails", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// ====== Read file storage ====== //
		emailStorageFile, err := os.Open("email.storage.json")

		if err != nil {
			fmt.Println(err)
		}
		defer emailStorageFile.Close()

		byteEmailStorageFile, _ := ioutil.ReadAll(emailStorageFile)

		var emails EmailStorage
		if err := json.Unmarshal(byteEmailStorageFile, &emails); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}
		// ====== ====== ====== //

		BTCRate := getBTCValue(w)

		messageText := fmt.Sprintf("Current BTC to UAH rate is: %f", BTCRate)

		for _, email := range emails {
			fmt.Println("emails:", email)
			sendBTCEmail(email, messageText)
		}

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

func arrayContains(arr []string, e string) bool {
	for _, a := range arr {
		if a == e {
			return true
		}
	}
	return false
}

func sendBTCEmail(email string, messageText string) {
	from := mail.NewEmail("Dmytro Uchkin", "dimon97xl+sandgrid@gmail.com") // Change to your verified sender

	subject := "BTC to UAH Rate"

	to := mail.NewEmail("Recipient", email) // Change to your recipient

	htmlContent := "<span>" + messageText + "</span>"
	message := mail.NewSingleEmail(from, subject, to, messageText, htmlContent)
	client := sendgrid.NewSendClient("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")
	response, err := client.Send(message)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Headers)
	}
}

func getBTCValue(w http.ResponseWriter) float64 {
	resp, err := http.Get("http://api.coinlayer.com/live?access_key=06f05f91a0c78ceb874adc4d6e65bdb2&target=UAH")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var result CoinlayerResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	return result.Rates.BTC
}
