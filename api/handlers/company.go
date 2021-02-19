package handlers

import "net/http"

// CompanyHandler ...
type CompanyHandler struct{}

// NewCompanyHandler ...
func NewCompanyHandler() *CompanyHandler {
	return &CompanyHandler{}
}

func (e *CompanyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Write([]byte("GET"))
	case http.MethodPost:
		w.Write([]byte("POST"))
	case http.MethodPut:
		w.Write([]byte("PUT"))
	case http.MethodDelete:
		w.Write([]byte("DELETE"))
	}
}
