package servers

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/Felley/accounting-service/protos/accounting"
)

// EmployeeServer ...
type EmployeeServer struct {
	db *sql.DB
	l  *log.Logger
	mu *sync.Mutex
	accounting.UnimplementedEmployeeAccountingServer
}

// NewEmployeeServer returns new employee storage processing server
func NewEmployeeServer(db *sql.DB, l *log.Logger, mu *sync.Mutex) *EmployeeServer {
	return &EmployeeServer{db: db, l: l, mu: mu}
}

// AddEmployee adds employee to DB
func (es *EmployeeServer) AddEmployee(ctx context.Context, req *accounting.EmployeeRequest) (*accounting.EmployeeResponce, error) {
	var query string
	if req.ID != 0 {
		return es.UpdateEmployee(ctx, req)
	}

	query = fmt.Sprintf("INSERT employee (name, second_name, surname, hire_date, position, company_id) VALUES ('%s', '%s', '%s', '%s', '%s', %d);",
		req.Name, req.SecondName, req.Surname, req.HireDate, req.Position, req.CompanyID)

	es.mu.Lock()
	_, err := es.db.Exec(query)
	es.mu.Unlock()
	if err != nil {
		es.l.Printf("%s occured while executing AddEmployee SQL query", err.Error())
		return nil, err
	}
	return &accounting.EmployeeResponce{StatusCode: 200}, nil
}

// UpdateEmployee updates employee specified info
func (es *EmployeeServer) UpdateEmployee(ctx context.Context, req *accounting.EmployeeRequest) (*accounting.EmployeeResponce, error) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("UPDATE employee SET name = '%s'", req.Name))
	if req.SecondName != "" {
		buffer.WriteString(fmt.Sprintf(", second_name = '%s'", req.SecondName))
	}
	if req.Surname != "" {
		buffer.WriteString(fmt.Sprintf(", surname = '%s'", req.Surname))
	}
	if req.HireDate != "" {
		buffer.WriteString(fmt.Sprintf(", hire_date = '%s'", req.HireDate))
	}
	if req.Position != "" {
		buffer.WriteString(fmt.Sprintf(", position = '%s'", req.Position))
	}
	if req.CompanyID != 0 {
		buffer.WriteString(fmt.Sprintf(", company_id = %d", req.CompanyID))
	}
	buffer.WriteString(fmt.Sprintf(" WHERE id = %d", req.ID))
	query := buffer.String()

	es.mu.Lock()
	rows, err := es.db.Query(fmt.Sprintf("SELECT * FROM employee WHERE id = %d", req.ID))
	if err != nil {
		es.l.Printf("%s occured while executing UpdateEmployee SQL query", err.Error())
		return nil, err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		e := &accounting.EmployeeResponce{}
		err := rows.Scan(&e.ID, &e.Name, &e.SecondName, &e.Surname, &e.HireDate, &e.Position, &e.CompanyID)
		if err != nil {
			es.l.Println(err)
			return nil, err
		}
		found = true
	}
	if !found {
		return nil, errors.New("Employee not found")
	}
	_, _ = es.db.Exec(query)
	es.mu.Unlock()
	return &accounting.EmployeeResponce{StatusCode: 200}, nil
}

// GetEmployee gets employee by specified id
func (es *EmployeeServer) GetEmployee(ctx context.Context, req *accounting.EmployeeRequest) (*accounting.EmployeeResponce, error) {
	query := fmt.Sprintf("SELECT * FROM employee WHERE id = %d", req.ID)

	es.mu.Lock()
	rows, err := es.db.Query(query)
	es.mu.Unlock()
	if err != nil {
		es.l.Printf("%s occured while executing UpdateEmployee SQL query", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := &accounting.EmployeeResponce{}
		err := rows.Scan(&e.ID, &e.Name, &e.SecondName, &e.Surname, &e.HireDate, &e.Position, &e.CompanyID)
		if err != nil {
			es.l.Println(err)
			return nil, err
		}
		e.StatusCode = 200
		return e, nil
	}
	return nil, errors.New("Employee not found")
}

// DeleteEmployee deletes employee by specified id
func (es *EmployeeServer) DeleteEmployee(ctx context.Context, req *accounting.EmployeeRequest) (*accounting.EmployeeResponce, error) {
	query := fmt.Sprintf("DELETE FROM employee WHERE id = %d", req.ID)

	es.mu.Lock()
	res, err := es.db.Exec(query)
	es.mu.Unlock()
	if err != nil {
		es.l.Printf("%s occured while executing UpdateEmployee SQL query", err.Error())
		return nil, err
	}

	if n, err := res.RowsAffected(); err != nil || n == 0 {
		return nil, errors.New("Employee not found")
	}
	return &accounting.EmployeeResponce{StatusCode: 200}, nil
}
