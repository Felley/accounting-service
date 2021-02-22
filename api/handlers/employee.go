package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
func (e *EmployeeHandler) AddEmployee(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	employee := &data.Employee{}
	err := employee.FromJSON(req.Body)
	if err != io.EOF {
		http.Error(w, "Invalid input", 405)
		return
	}

	r := &accounting.EmployeeRequest{
		ID:         employee.ID,
		Name:       employee.Name,
		SecondName: employee.SecondName,
		Surname:    employee.Surname,
		HireDate:   employee.HireDate,
		Position:   employee.Position,
		CompanyID:  employee.CompanyID,
	}

	_, err = e.ec.AddEmployee(context.Background(), r)
	if err != nil {
		http.Error(w, "Invalid input", 405)
		return
	}
}

// GetEmployee looks for employee by specified id
func (e *EmployeeHandler) GetEmployee(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID supplied", 400)
		return
	}
	resp, err := e.ec.GetEmployee(context.Background(), &accounting.EmployeeRequest{ID: id})
	if err != nil {
		http.Error(w, "Employee not found", 404)
		return
	}
	enc := json.NewEncoder(w)
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
		http.Error(w, "Unexpected error while sending answer", 404)
		return
	}
	fmt.Println(id)
}

// PostFormEmployee updates employee data by incomming form data
func (e *EmployeeHandler) PostFormEmployee(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := req.ParseMultipartForm(128 * 1024)
	if err != nil {
		http.Error(w, "Invalid input", 405)
		w.WriteHeader(405)
		return
	}
}

// UpdateEmployee updates employee data by incomming json data
func (e *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	_, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID supplied", 400)
		return
	}

	employee := &data.Employee{}
	err = employee.FromJSON(req.Body)
	if err != nil {
		http.Error(w, "Invalid input", 405)
		return
	}
	r := &accounting.EmployeeRequest{
		ID:         employee.ID,
		Name:       employee.Name,
		SecondName: employee.SecondName,
		Surname:    employee.Surname,
		HireDate:   employee.HireDate,
		Position:   employee.Position,
		CompanyID:  employee.CompanyID,
	}
	_, err = e.ec.UpdateEmployee(context.Background(), r)

	if err != nil {
		http.Error(w, "Unexpected error while sending answer", 404)
	}
}

// DeleteEmployee deletes employee by specified id
func (e *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "	Invalid ID supplied", 400)
		return
	}
	_, err = e.ec.DeleteEmployee(context.Background(), &accounting.EmployeeRequest{ID: id})
	if err != nil {
		http.Error(w, "Unexpected error while sending answer", 404)
	}
}

// MiddlewareEmployeeValidation validates incoming json data
func (e *EmployeeHandler) MiddlewareEmployeeValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		employee := &data.Employee{}

		err := employee.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Invalid input", 405)
			return
		}

		err = employee.Validate()
		if err != nil {
			http.Error(rw, "Invalid input", 405)
			return
		}
		ctx := context.WithValue(r.Context(), data.Employee{}, employee)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
