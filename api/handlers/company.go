package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Felley/accounting-service/api/data"
	"github.com/Felley/accounting-service/protos/accounting"
	"github.com/gorilla/mux"
)

// CompanyHandler ...
type CompanyHandler struct {
	l  *log.Logger
	cc accounting.CompanyAccountingClient
}

// NewCompanyHandler ...
func NewCompanyHandler(l *log.Logger, cc accounting.CompanyAccountingClient) *CompanyHandler {
	return &CompanyHandler{l, cc}
}

// AddCompany ...
func (c *CompanyHandler) AddCompany(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	value := r.Context().Value(companyKey{})
	company := value.(*data.Company)

	req := &accounting.CompanyRequest{
		ID:        company.ID,
		Name:      company.Name,
		LegalForm: company.LegalForm,
	}

	_, err := c.cc.AddCompany(context.Background(), req)
	if err != nil {
		http.Error(rw, "Invalid input", 405)
		return
	}
}

// UpdateCompany ...
func (c *CompanyHandler) UpdateCompany(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	value := r.Context().Value(companyKey{})
	company := value.(*data.Company)

	req := &accounting.CompanyRequest{
		ID:        company.ID,
		Name:      company.Name,
		LegalForm: company.LegalForm,
	}

	_, err := c.cc.UpdateCompany(context.Background(), req)
	if err != nil {
		http.Error(rw, "Invalid input", 405)
		return
	}
}

// GetCompany ...
func (c *CompanyHandler) GetCompany(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(rw, "Invalid ID supplied", 400)
		return
	}
	resp, err := c.cc.GetCompany(context.Background(), &accounting.CompanyRequest{ID: id})
	if err != nil {
		http.Error(rw, "Employee not found", 404)
		return
	}
	enc := json.NewEncoder(rw)
	err = enc.Encode(&data.Company{
		ID:        resp.ID,
		Name:      resp.Name,
		LegalForm: resp.LegalForm,
	})
	if err != nil {
		http.Error(rw, "Unexpected error while sending answer", 404)
		return
	}
}

// PostFormCompany ...
func (c *CompanyHandler) PostFormCompany(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)

	if err != nil {
		http.Error(rw, "Invalid ID supplied", 400)
		return
	}
	err = r.ParseMultipartForm(128 * 1024)
	if err != nil {
		http.Error(rw, "Invalid input", 405)
		return
	}

	company := &data.Company{
		ID:        id,
		Name:      r.Form.Get("id"),
		LegalForm: r.Form.Get("legalForm"),
	}

	err = company.Validate()
	if err != nil {
		http.Error(rw, "Invalid input", 405)
		return
	}

	req := &accounting.CompanyRequest{
		ID:        company.ID,
		Name:      company.Name,
		LegalForm: company.LegalForm,
	}
	_, err = c.cc.UpdateCompany(context.Background(), req)

	if err != nil {
		http.Error(rw, "Unexpected error while sending answer", 404)
	}
}

// DeleteCompany ...
func (c *CompanyHandler) DeleteCompany(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(rw, "Invalid ID supplied", 400)
		return
	}
	_, err = c.cc.DeleteCompany(context.Background(), &accounting.CompanyRequest{ID: id})
	if err != nil {
		http.Error(rw, "Unexpected error while sending answer", 404)
	}
}

type companyKey struct{}

// MiddlewareCompanyValidation ...
func (c *CompanyHandler) MiddlewareCompanyValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		company := &data.Company{}

		err := company.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Invalid input", 405)
			return
		}
		defer r.Body.Close()
		err = company.Validate()
		if err != nil {
			http.Error(rw, "Invalid input", 405)
			return
		}
		ctx := context.WithValue(r.Context(), companyKey{}, company)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
