package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/Felley/accounting-service/api/data"
	"github.com/gorilla/mux"
)

// EmployeeHandler ...
type EmployeeHandler struct {
	l *log.Logger
}

// NewEmployeeHandler ...
func NewEmployeeHandler(l *log.Logger) *EmployeeHandler {
	return &EmployeeHandler{l}
}

// AddEmployee ...
func (e *EmployeeHandler) AddEmployee(w http.ResponseWriter, req *http.Request) {
	e.l.Println("Employee add person called")
	employee := &data.Employee{}
	err := employee.FromJSON(req.Body)
	if err != nil {
		http.Error(w, "Invalid input", 405)
	}

	e.l.Printf("Person: %#v", employee)
}

// ErrEmployeeNotFound ...
var ErrEmployeeNotFound = fmt.Errorf("Employee not found")

// findEmployee ...
func (e *EmployeeHandler) findEmployee(id int64) (*data.Employee, error) {
	var employeeList []*data.Employee
	var found bool
	for _, employee := range employeeList {
		if employee.ID == id {
			found = true
		}
	}
	if !found {
		return nil, ErrEmployeeNotFound
	}
	return nil, nil
}

// GetEmployee ...
func (e *EmployeeHandler) GetEmployee(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(400)
	}
	fmt.Println(id)
}

// UpdateEmployee ...
func (e *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, req *http.Request) {
	// FIXME: fix bug with links ends by /1223ssasdad
	r := regexp.MustCompile(`/([0-9]+)\/{0,1}`)
	g := r.FindAllStringSubmatch(req.URL.Path, -1)
	if len(g) != 1 {
		http.Error(w, "Invalid ID supplied", 400)
		return
	}

	if len(g[0]) != 2 {
		http.Error(w, "Invalid ID supplied", 400)
		return
	}

	id, err := strconv.Atoi(g[0][1])
	if err != nil {
		http.Error(w, "Invalid ID supplied", 400)
		return
	}
	e.l.Println("got id", id)
}

// DeleteEmployee ...
func (e *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, req *http.Request) {
}

// KeyEmployee ...
type KeyEmployee struct{}

// MiddlewareEmployeeValidation ...
func (e *EmployeeHandler) MiddlewareEmployeeValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		employee := &data.Employee{}

		err := employee.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Invalid input", 405)
		}

		ctx := context.WithValue(r.Context(), KeyEmployee{}, employee)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
