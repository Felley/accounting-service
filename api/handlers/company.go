package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Felley/accounting-service/api/data"
	"github.com/Felley/accounting-service/protos/accounting"
)

// CompanyHandler struct is a handler struct which is gRPC company client,
// it also logs errors in terminal
type CompanyHandler struct {
	l  *log.Logger
	cc accounting.CompanyAccountingClient
}

// NewCompanyHandler creates handler for company API processing
func NewCompanyHandler(l *log.Logger, cc accounting.CompanyAccountingClient) *CompanyHandler {
	return &CompanyHandler{l, cc}
}

// NewCompanyRequest creates CompanyRequestStruct filling it's data
func NewCompanyRequest(id int64, name string, legalForm string) *accounting.CompanyRequest {
	return &accounting.CompanyRequest{
		ID:        id,
		Name:      name,
		LegalForm: legalForm,
	}
}

// AddCompany sends query for adding company to DB
func (c *CompanyHandler) AddCompany(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	value := r.Context().Value(companyKey{})
	company := value.(*data.Company)

	req := NewCompanyRequest(company.ID, company.Name, company.LegalForm)

	_, err := c.cc.AddCompany(context.Background(), req)
	if err != nil {
		c.l.Printf("%e ocuured while adding company to DB", err)
		http.Error(rw, "Invalid input", 405)
		return
	}
}

// UpdateCompany updates company data in DB
func (c *CompanyHandler) UpdateCompany(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	value := r.Context().Value(companyKey{})
	company := value.(*data.Company)

	req := NewCompanyRequest(company.ID, company.Name, company.LegalForm)

	_, err := c.cc.UpdateCompany(context.Background(), req)
	if err != nil {
		c.l.Printf("%e ocuured while updating company info", err)
		http.Error(rw, "Employee not found", 404)
		return
	}
}

// GetCompany looks for company in DB by specified id
func (c *CompanyHandler) GetCompany(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	id, err := extractIDFromLink(rw, r)
	if err != nil || id < 1 {
		http.Error(rw, "Invalid ID supplied", 400)
		return
	}

	resp, err := c.cc.GetCompany(context.Background(), &accounting.CompanyRequest{ID: id})
	if err != nil {
		http.Error(rw, "Company not found", 404)
		return
	}

	enc := json.NewEncoder(rw)
	err = enc.Encode(NewCompanyRequest(resp.ID, resp.Name, resp.LegalForm))
	if err != nil {
		c.l.Printf("%e ocuured while sending answer", err)
		http.Error(rw, "Unexpected error while sending answer", 404)
		return
	}
}

// PostFormCompany updates company data by incomming form data
func (c *CompanyHandler) PostFormCompany(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	id, err := extractIDFromLink(rw, r)
	if err != nil {
		http.Error(rw, "Invalid input", 400)
		return
	}
	err = r.ParseMultipartForm(128 * 1024)
	if err != nil {
		http.Error(rw, "Invalid input", 405)
		return
	}

	company := data.NewCompany(id, r.Form.Get("name"), r.Form.Get("status"))

	err = company.Validate()
	if err != nil {
		http.Error(rw, "Invalid input", 405)
		return
	}

	req := NewCompanyRequest(company.ID, company.Name, company.LegalForm)
	_, err = c.cc.UpdateCompany(context.Background(), req)

	if err != nil {
		c.l.Printf("%e ocuured while updating company info", err)
		http.Error(rw, "Unexpected error while sending answer", 404)
	}
}

// DeleteCompany deletes company by specified id
func (c *CompanyHandler) DeleteCompany(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	id, err := extractIDFromLink(rw, r)
	if err != nil {
		http.Error(rw, "Invalid ID supplied", 400)
		return
	}

	_, err = c.cc.DeleteCompany(context.Background(), &accounting.CompanyRequest{ID: id})
	if err != nil {
		c.l.Printf("%e ocuured while sending answer", err)
		http.Error(rw, "Unexpected error while sending answer", 404)
	}
}

// GetCompanyEmployees get's company employees
func (c *CompanyHandler) GetCompanyEmployees(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	id, err := extractIDFromLink(rw, r)
	if err != nil {
		http.Error(rw, "Invalid ID supplied", 400)
		return
	}

	resp, err := c.cc.GetCompanyEmployees(context.Background(), &accounting.CompanyRequest{ID: id})
	if err != nil {
		http.Error(rw, "Company not found", 404)
		return
	}

	enc := json.NewEncoder(rw)
	err = enc.Encode(resp.Employees)
	if err != nil {
		c.l.Printf("%e ocuured while sending answer", err)
		http.Error(rw, "Unexpected error while sending answer", 404)
		return
	}
}

type companyKey struct{}

// MiddlewareCompanyValidation validates incoming json data
func (c *CompanyHandler) MiddlewareCompanyValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		company := &data.Company{}

		err := company.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Validation exception", 405)
			return
		}
		defer r.Body.Close()

		if company.ID < 1 {
			http.Error(rw, "Invalid ID supplied", 400)
			return
		}

		err = company.Validate()
		if err != nil {
			http.Error(rw, "Validation exception", 405)
			return
		}
		ctx := context.WithValue(r.Context(), companyKey{}, company)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
