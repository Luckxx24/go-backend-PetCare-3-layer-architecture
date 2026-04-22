package jsonresponse

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("errors %s", msg)
	}
	type Errors struct {
		Error string `json:"error"`
	}
}

func ResponSuccess(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		RespondError(w, 500, "gagal marshall dari database")
	}

	w.Header().Add("contetn-type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

}
