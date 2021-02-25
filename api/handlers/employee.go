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

// EmployeeHandler struct is a handler struct which is gRPC employee client,
// it also logs errors in terminal
type EmployeeHandler struct {
	l  *log.Logger
	ec accounting.EmployeeAccountingClient
}

// NewEmployeeHandler creates handler for employee API processing
func NewEmployeeHandler(l *log.Logger, ec accounting.EmployeeAccountingClient) *EmployeeHandler {
	return &EmployeeHandler{l, ec}
}

// NewEmployeeRequest creates EmployeeRequestStruct filling it's data
func NewEmployeeRequest(id int64, name string, secondName string, surname string, hireDate string, position string, companyID int64) *accounting.EmployeeRequest {
	return &accounting.EmployeeRequest{
		ID:         id,
		Name:       name,
		SecondName: secondName,
		Surname:    surname,
		HireDate:   hireDate,
		Position:   position,
		CompanyID:  companyID,
	}
}

func extractIDFromLink(rw http.ResponseWriter, r *http.Request) (id int64, err error) {
	vars := mux.Vars(r)
	id, err = strconv.ParseInt(vars["id"], 10, 64)
	return id, err
}

// AddEmployee sends query for adding employee to DB
func (e *EmployeeHandler) AddEmployee(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	value := r.Context().Value(employeeKey{})
	employee := value.(*data.Employee)

	req := NewEmployeeRequest(employee.ID, employee.Name, employee.SecondName, employee.Surname, employee.HireDate, employee.Position, employee.CompanyID)

	_, err := e.ec.AddEmployee(context.Background(), req)
	if err != nil {
		e.l.Printf("%e ocuured while adding employee to DB", err)
		http.Error(rw, "Invalid input", 405)
		return
	}
}

// UpdateEmployee updates employee data by incomming json data
func (e *EmployeeHandler) UpdateEmployee(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	value := r.Context().Value(employeeKey{})
	employee := value.(*data.Employee)

	req := NewEmployeeRequest(employee.ID, employee.Name, employee.SecondName, employee.Surname, employee.HireDate, employee.Position, employee.CompanyID)
	_, err := e.ec.UpdateEmployee(context.Background(), req)
	if err != nil {
		e.l.Printf("%e ocuured while sending answer", err)
		http.Error(rw, "Unexpected error while sending answer", 404)
	}
}

// GetEmployee looks for employee in DB by specified id
func (e *EmployeeHandler) GetEmployee(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	id, err := extractIDFromLink(rw, r)
	if err != nil || id < 1 {
		http.Error(rw, "Invalid ID supplied", 400)
		return
	}

	resp, err := e.ec.GetEmployee(context.Background(), &accounting.EmployeeRequest{ID: id})
	if err != nil {
		http.Error(rw, "Employee not found", 404)
		return
	}

	enc := json.NewEncoder(rw)
	err = enc.Encode(NewEmployeeRequest(resp.ID, resp.Name, resp.SecondName, resp.Surname, resp.HireDate, resp.Position, resp.CompanyID))
	if err != nil {
		e.l.Printf("%e ocuured while sending answer", err)
		http.Error(rw, "Unexpected error while sending answer", 404)
		return
	}
}

// PostFormEmployee updates employee data by incomming form data
func (e *EmployeeHandler) PostFormEmployee(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	id, err := extractIDFromLink(rw, r)
	if err != nil {
		http.Error(rw, "Invalid input", 400)
		return
	}

	err = r.ParseMultipartForm(128 * 1024)
	if err != nil {
		e.l.Printf("%s while parsing form data", err.Error())
		http.Error(rw, "Invalid input", 405)
		return
	}

	companyID, err := strconv.ParseInt(r.Form.Get("companyId"), 10, 64)
	if err != nil {
		e.l.Printf("%s while parsing id", err.Error())
		http.Error(rw, "Invalid input", 405)
		return
	}

	employee := data.NewEmployee(id, r.Form.Get("name"), r.Form.Get("secondName"), r.Form.Get("surname"), r.Form.Get("hireDate"), r.Form.Get("position"), companyID)
	err = employee.Validate()
	if err != nil {
		e.l.Printf("%s while validating table data", err.Error())
		http.Error(rw, "Invalid input", 405)
		return
	}

	req := NewEmployeeRequest(employee.ID, employee.Name, employee.SecondName, employee.Surname, employee.HireDate, employee.Position, employee.CompanyID)
	_, err = e.ec.UpdateEmployee(context.Background(), req)
	if err != nil {
		e.l.Printf("%e ocuured while sending answer", err)
		http.Error(rw, "Unexpected error while sending answer", 404)
	}

}

// DeleteEmployee deletes employee by specified id
func (e *EmployeeHandler) DeleteEmployee(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	id, err := extractIDFromLink(rw, r)
	if err != nil || id < 1 {
		http.Error(rw, "Invalid ID supplied", 400)
		return
	}

	_, err = e.ec.DeleteEmployee(context.Background(), &accounting.EmployeeRequest{ID: id})
	if err != nil {
		http.Error(rw, "Employee not found", 404)
		return
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
			http.Error(rw, "Validation exception", 405)
			return
		}
		defer r.Body.Close()

		if employee.ID < 1 {
			http.Error(rw, "Invalid ID supplied", 400)
			return
		}

		err = employee.Validate()
		if err != nil {
			http.Error(rw, "Validation exception", 405)
			return
		}

		ctx := context.WithValue(r.Context(), employeeKey{}, employee)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}
