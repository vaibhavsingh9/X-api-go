package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
		fmt.Println("error loading .env file")
	}
	fmt.Println("starting server")

	if args := os.Args; len(args) > 1 && args[1] == "-register" {
		go registerWebhook()
	}

	m := mux.NewRouter()
	m.HandleFunc("/", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(200)
		fmt.Fprintf(writer, "server is up and running")
	})
	m.HandleFunc("/webhook/twitter", CrcCheck).Methods("GET")

	server := &http.Server{
		Handler: m,
	}
	server.Addr = ":9090"
	server.ListenAndServe()
}

func CrcCheck(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	token := request.URL.Query()["crc_token"]
	if len(token) < 1 {
		fmt.Fprintf(writer, "no crc_token given")
		return
	}
	h := hmac.New(sha256.New, []byte(os.Getenv("CONSUMER_SECRET")))
	h.Write([]byte(token[0]))
	encoded := base64.StdEncoding.EncodeToString(h.Sum(nil))
	response := make(map[string]string)
	response["response_token"] = "sha256" + encoded
	responseJson, _ := json.Marshal(response)
	fmt.Fprintf(writer, string(responseJson))
}
