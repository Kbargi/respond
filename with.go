package respond

import (
	"encoding/json"
	"log"
	"net/http"
)

// With specifies the status code and data to repsond with.
func With(status int, data interface{}) *W {
	return &W{Code: status, Data: data}
}

// W holds details about the response that will be made
// when To is called.
type W struct {
	Code   int
	Data   interface{}
	header http.Header
}

// To writes the repsonse.
func (with *W) To(w http.ResponseWriter, r *http.Request) {
	// copy headers
	copyheaders(with.header, w.Header())
	// find the encoder
	encoder, ok := Encoders().Match(r.Header.Get("Accept"))
	if !ok {
		encoder = DefaultEncoder
	}
	// write response
	if err := Write(w, r, with.Code, with.Data, encoder); err != nil {
		Err(w, r, with, err)
	}
}

// Write is the function that actually writes the response.
var Write = func(w http.ResponseWriter, r *http.Request, status int, data interface{}, encoder Encoder) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// Err is called when an internal error occurs while responding.
var Err = func(w http.ResponseWriter, r *http.Request, with *W, err error) {
	log.Println()
}
