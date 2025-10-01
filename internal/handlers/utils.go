package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func sendHttpMessage[T any](w http.ResponseWriter, data T, code int) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("error while marshaling message data:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonData)

	log.Printf("Sent HTTP message '%s' with code %d\n", jsonData, code)
}

func parseHttpJson[T any](jsonData io.ReadCloser) (T, error) {
	var obj T
	if err := json.NewDecoder(jsonData).Decode(&obj); err != nil {
		return obj, err
	}

	return obj, nil
}
