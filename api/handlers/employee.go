package handlers

import (
	"context"
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
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID supplied", 400)
	}
	fmt.Println(id)
}

// PostFormEmployee updates employee data by incomming form data
func (e *EmployeeHandler) PostFormEmployee(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(128 * 1024)
	if err != nil {
		http.Error(w, "Invalid input", 405)
		w.WriteHeader(405)
		return
	}
}

// UpdateEmployee updates employee data by incomming json data
func (e *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "	Invalid ID supplied", 400)
	}

	employee := &data.Employee{}
	err = employee.FromJSON(req.Body)
	if err != nil {
		http.Error(w, "Invalid input", 405)
	}
	e.l.Println("got id", id)
}

// DeleteEmployee deletes employee by specified id
func (e *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, req *http.Request) {
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
