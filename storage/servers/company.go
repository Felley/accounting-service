package servers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Felley/accounting-service/protos/accounting"
)

// CompanyServer ...
type CompanyServer struct {
	db *sql.DB
	l  *log.Logger
	accounting.UnimplementedCompanyAccountingServer
}

// NewCompanyServer returns new employee storage processing server
func NewCompanyServer(db *sql.DB, l *log.Logger) *CompanyServer {
	return &CompanyServer{db: db, l: l}
}

// AddCompany ...
func (cs *CompanyServer) AddCompany(ctx context.Context, req *accounting.CompanyRequest) (*accounting.CompanyResponce, error) {
	var query string
	if req.ID != 0 {
		return cs.UpdateCompany(ctx, req)
	}

	query = fmt.Sprintf("INSERT company (name, legal_form) VALUES ('%s', '%s');", req.Name, req.LegalForm)

	fmt.Println(query)
	_, err := cs.db.Query(query)
	if err != nil {
		cs.l.Printf("%s occured while executing AddCompany SQL query", err.Error())
		return nil, err
	}

	return &accounting.CompanyResponce{StatusCode: 200}, nil
}

// UpdateCompany ...
func (cs *CompanyServer) UpdateCompany(ctx context.Context, req *accounting.CompanyRequest) (*accounting.CompanyResponce, error) {
	query := fmt.Sprintf("UPDATE company SET name = '%s', legal_form = '%s' WHERE id = %d", req.Name, req.LegalForm, req.ID)

	fmt.Println(query)
	_, err := cs.db.Query(query)
	if err != nil {
		cs.l.Printf("%s occured while executing UpdateCompany SQL query", err.Error())
		return nil, err
	}
	return &accounting.CompanyResponce{StatusCode: 200}, nil
}

// GetCompany ...
func (cs *CompanyServer) GetCompany(ctx context.Context, req *accounting.CompanyRequest) (*accounting.CompanyResponce, error) {
	query := fmt.Sprintf("SELECT * FROM company WHERE id = %d", req.ID)

	fmt.Println(query)
	rows, err := cs.db.Query(query)
	if err != nil {
		cs.l.Printf("%s occured while executing GetCompany SQL query", err.Error())
		return nil, err
	}

	for rows.Next() {
		e := &accounting.CompanyResponce{}
		err := rows.Scan(&e.ID, &e.Name, &e.LegalForm)
		if err != nil {
			cs.l.Println(err)
			return nil, err
		}
		e.StatusCode = 200
		return e, nil
	}

	return nil, errors.New("Employee not found")
}

// GetCompanyEmployees ...
func (cs *CompanyServer) GetCompanyEmployees(ctx context.Context, req *accounting.CompanyRequest) (*accounting.CompanyEmployeesResponce, error) {
	query := fmt.Sprintf("SELECT * FROM employee WHERE company_id = %d", req.ID)

	fmt.Println(query)
	rows, err := cs.db.Query(query)
	if err != nil {
		cs.l.Printf("%s occured while executing GetCompany SQL query", err.Error())
		return nil, err
	}

	var employees []*accounting.EmployeeResponce
	counter := 0
	for rows.Next() {
		counter++
		e := &accounting.EmployeeResponce{}
		err := rows.Scan(&e.ID, &e.Name, &e.SecondName, &e.Surname, &e.HireDate, &e.Position, &e.CompanyID)
		if err != nil {
			cs.l.Println(err)
			return nil, err
		}
		employees = append(employees, e)
	}

	if counter == 0 {
		return nil, errors.New("Company not found")
	}
	return &accounting.CompanyEmployeesResponce{Employees: employees}, nil
}

// DeleteCompany ...
func (cs *CompanyServer) DeleteCompany(ctx context.Context, req *accounting.CompanyRequest) (*accounting.CompanyResponce, error) {
	query := fmt.Sprintf("DELETE FROM company WHERE id = %d", req.ID)

	fmt.Println(query)
	_, err := cs.db.Query(query)
	if err != nil {
		cs.l.Printf("%s occured while executing DeleteCompany SQL query", err.Error())
		return nil, err
	}
	return &accounting.CompanyResponce{StatusCode: 200}, nil
}
