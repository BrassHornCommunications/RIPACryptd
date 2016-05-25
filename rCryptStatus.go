package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"net/http"
	"strconv"
)

func rCryptStatus(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	w.Header().Set("Content-Type", "application/json")
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	response := NewAPIResponse()
	response.Message = "All is good!"
	js, err := json.Marshal(response)
	if err == nil {
		w.WriteHeader(response.StatusCode)
		w.Write(js)
	} else {
		http.Error(w, err.Error(), 500)
	}

}
