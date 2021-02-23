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

// EmployeeHandler handles ...
type EmployeeHandler struct {
	l  *log.Logger
	ec accounting.EmployeeAccountingClient
}

// NewEmployeeHandler creates handler for employee API processing
func NewEmployeeHandler(l *log.Logger, ec accounting.EmployeeAccountingClient) *EmployeeHandler {
	return &EmployeeHandler{l, ec}
}

// AddEmployee sends query for adding employee to DB
func (e *EmployeeHandler) AddEmployee(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	value := r.Context().Value(employeeKey{})
	employee := value.(*data.Employee)

	req := &accounting.EmployeeRequest{
		ID:         employee.ID,
		Name:       employee.Name,
		SecondName: employee.SecondName,
		Surname:    employee.Surname,
		HireDate:   employee.HireDate,
		Position:   employee.Position,
		CompanyID:  employee.CompanyID,
	}

	_, err := e.ec.AddEmployee(context.Background(), req)
	if err != nil {
		http.Error(rw, "Invalid input", 405)
		return
	}
}

// GetEmployee looks for employee by specified id
func (e *EmployeeHandler) GetEmployee(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(rw, "Invalid ID supplied", 400)
		return
	}
	resp, err := e.ec.GetEmployee(context.Background(), &accounting.EmployeeRequest{ID: id})
	if err != nil {
		http.Error(rw, "Employee not found", 404)
		return
	}
	enc := json.NewEncoder(rw)
	err = enc.Encode(&data.Employee{
		ID:         resp.ID,
		Name:       resp.Name,
		SecondName: resp.SecondName,
		Surname:    resp.Surname,
		HireDate:   resp.HireDate,
		Position:   resp.Position,
		CompanyID:  resp.CompanyID,
	})
	if err != nil {
		http.Error(rw, "Unexpected error while sending answer", 404)
		return
	}
}

// PostFormEmployee updates employee data by incomming form data
func (e *EmployeeHandler) PostFormEmployee(rw http.ResponseWriter, r *http.Request) {
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
	companyID, err := strconv.ParseInt(r.Form.Get("companyId"), 10, 64)
	if err != nil {
		http.Error(rw, "Invalid input", 405)
		return
	}
	employee := &data.Employee{
		ID:         id,
		Name:       r.Form.Get("id"),
		SecondName: r.Form.Get("id"),
		Surname:    r.Form.Get("id"),
		HireDate:   r.Form.Get("id"),
		Position:   r.Form.Get("id"),
		CompanyID:  companyID,
	}
	err = employee.Validate()
	if err != nil {
		http.Error(rw, "Invalid input", 405)
		return
	}
	req := &accounting.EmployeeRequest{
		ID:         employee.ID,
		Name:       employee.Name,
		SecondName: employee.SecondName,
		Surname:    employee.Surname,
		HireDate:   employee.HireDate,
		Position:   employee.Position,
		CompanyID:  employee.CompanyID,
	}
	_, err = e.ec.UpdateEmployee(context.Background(), req)

	if err != nil {
		http.Error(rw, "Unexpected error while sending answer", 404)
	}

}

// UpdateEmployee updates employee data by incomming json data
func (e *EmployeeHandler) UpdateEmployee(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	_, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(rw, "Invalid ID supplied", 400)
		return
	}
	value := r.Context().Value(employeeKey{})
	employee := value.(*data.Employee)

	req := &accounting.EmployeeRequest{
		ID:         employee.ID,
		Name:       employee.Name,
		SecondName: employee.SecondName,
		Surname:    employee.Surname,
		HireDate:   employee.HireDate,
		Position:   employee.Position,
		CompanyID:  employee.CompanyID,
	}
	_, err = e.ec.UpdateEmployee(context.Background(), req)

	if err != nil {
		http.Error(rw, "Unexpected error while sending answer", 404)
	}
}

// DeleteEmployee deletes employee by specified id
func (e *EmployeeHandler) DeleteEmployee(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(rw, "Invalid ID supplied", 400)
		return
	}
	_, err = e.ec.DeleteEmployee(context.Background(), &accounting.EmployeeRequest{ID: id})
	if err != nil {
		http.Error(rw, "Unexpected error while sending answer", 404)
	}
}

type employeeKey struct{}

// MiddlewareEmployeeValidation validates incoming json data
func (e *EmployeeHandler) MiddlewareEmployeeValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		employee := &data.Employee{}

		err := employee.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Invalid input", 405)
			return
		}
		defer r.Body.Close()
		err = employee.Validate()
		if err != nil {
			http.Error(rw, "Invalid input", 405)
			return
		}
		ctx := context.WithValue(r.Context(), employeeKey{}, employee)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
