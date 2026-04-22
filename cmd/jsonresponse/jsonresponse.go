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

	Writejson(w, code, Errors{
		Error: msg,
	})
}

func ResponSuccess(w http.ResponseWriter, code int, payload interface{}) {
	err := Writejson(w, code, payload)

	if err != nil {
		log.Printf("[MARSHALL ERROR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error 500"}`))
	}

}

func Writejson(w http.ResponseWriter, code int, payload interface{}) error {
	data, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

	return nil

}

func RespondWithBadRequest(w http.ResponseWriter, msg string) {
	RespondError(w, 400, msg)
}

func RespondWithUnauthorized(w http.ResponseWriter, msg string) {
	RespondError(w, 401, msg)
}

func RespondWithNotfound(w http.ResponseWriter, msg string) {
	RespondError(w, 404, msg)
}

func RespondWithForbiden(w http.ResponseWriter, msg string) {
	RespondError(w, 403, msg)
}

func RespondWithConflict(w http.ResponseWriter, msg string) {
	RespondError(w, 409, msg)
}
