package servers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Felley/accounting-service/protos/accounting"
)

// EmployeeServer ...
type EmployeeServer struct {
	db *sql.DB
	l  *log.Logger
	accounting.UnimplementedEmployeeAccountingServer
}

// NewEmployeeServer returns new employee storage processing server
func NewEmployeeServer(db *sql.DB, l *log.Logger) *EmployeeServer {
	return &EmployeeServer{db: db, l: l}
}

// AddEmployee ...
func (es *EmployeeServer) AddEmployee(ctx context.Context, req *accounting.EmployeeRequest) (*accounting.EmployeeResponce, error) {
	var query string
	if req.ID != 0 {
		return es.UpdateEmployee(ctx, req)
	}

	query = fmt.Sprintf("INSERT employee (name, second_name, surname, hire_date, position, company_id) VALUES ('%s', '%s', '%s', '%s', '%s', %d);",
		req.Name, req.SecondName, req.Surname, req.HireDate, req.Position, req.CompanyID)

	fmt.Println(query)
	_, err := es.db.Query(query)
	if err != nil {
		es.l.Printf("%s occured while executing AddEmployee SQL query", err.Error())
		return nil, err
	}

	return &accounting.EmployeeResponce{StatusCode: 200}, nil
}

// UpdateEmployee ...
func (es *EmployeeServer) UpdateEmployee(ctx context.Context, req *accounting.EmployeeRequest) (*accounting.EmployeeResponce, error) {
	query := fmt.Sprintf("UPDATE employee SET name = '%s', second_name = '%s', surname = '%s', hire_date = '%s', position = '%s', company_id = %d WHERE id = %d",
		req.Name, req.SecondName, req.Surname, req.HireDate, req.Position, req.CompanyID, req.ID)

	fmt.Println(query)
	_, err := es.db.Query(query)
	if err != nil {
		es.l.Printf("%s occured while executing UpdateEmployee SQL query", err.Error())
		return nil, err
	}
	return &accounting.EmployeeResponce{StatusCode: 200}, nil
}

// GetEmployee ...
func (es *EmployeeServer) GetEmployee(ctx context.Context, req *accounting.EmployeeRequest) (*accounting.EmployeeResponce, error) {
	query := fmt.Sprintf("SELECT * FROM employee WHERE id = %d", req.ID)

	fmt.Println(query)
	rows, err := es.db.Query(query)
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

// DeleteEmployee ...
func (es *EmployeeServer) DeleteEmployee(ctx context.Context, req *accounting.EmployeeRequest) (*accounting.EmployeeResponce, error) {
	query := fmt.Sprintf("DELETE FROM employee WHERE id = %d", req.ID)

	fmt.Println(query)
	_, err := es.db.Query(query)
	if err != nil {
		es.l.Printf("%s occured while executing UpdateEmployee SQL query", err.Error())
		return nil, err
	}
	return &accounting.EmployeeResponce{StatusCode: 200}, nil
}
