package utils

import (
	"encoding/json"
	"net/http"

	"github.com/adeo/go-api-skeleton/storage/model"
)

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set(HeaderNameContentType, HeaderValueApplicationJSONUTF8)
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func JSONError(w http.ResponseWriter, e model.APIError) {
	if e.Headers != nil {
		for k, headers := range e.Headers {
			for _, headerValue := range headers {
				w.Header().Add(k, headerValue)
			}
		}
	}
	JSON(w, e.HTTPCode, e)
}

func JSONErrorWithMessage(w http.ResponseWriter, e model.APIError, message string) {
	e.Description = message
	JSONError(w, e)
}

func YAML(w http.ResponseWriter, status int, data string) {
	w.Header().Set(HeaderNameContentType, HeaderValueApplicationYAML)
	w.WriteHeader(status)
	w.Write([]byte(data))
}
