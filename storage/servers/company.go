package servers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/Felley/accounting-service/protos/accounting"
)

// CompanyServer struct is gRPC server, that processes company data updates
type CompanyServer struct {
	db *sql.DB
	l  *log.Logger
	mu *sync.Mutex
	accounting.UnimplementedCompanyAccountingServer
}

// NewCompanyServer returns new employee storage processing server
func NewCompanyServer(db *sql.DB, l *log.Logger, mu *sync.Mutex) *CompanyServer {
	return &CompanyServer{db: db, l: l, mu: mu}
}

// AddCompany adds company to DB
func (cs *CompanyServer) AddCompany(ctx context.Context, req *accounting.CompanyRequest) (*accounting.CompanyResponce, error) {
	var query string
	if req.ID != 0 {
		return cs.UpdateCompany(ctx, req)
	}

	query = fmt.Sprintf("INSERT company (name, legal_form) VALUES ('%s', '%s');", req.Name, req.LegalForm)

	cs.mu.Lock()
	_, err := cs.db.Exec(query)
	cs.mu.Unlock()
	if err != nil {
		cs.l.Printf("%s occured while executing AddCompany SQL query", err.Error())
		return nil, err
	}

	return &accounting.CompanyResponce{StatusCode: 200}, nil
}

// UpdateCompany updated company specified info
func (cs *CompanyServer) UpdateCompany(ctx context.Context, req *accounting.CompanyRequest) (*accounting.CompanyResponce, error) {
	query := fmt.Sprintf("UPDATE company SET name = '%s', legal_form = '%s' WHERE id = %d", req.Name, req.LegalForm, req.ID)
	if req.LegalForm == "" {
		query = fmt.Sprintf("UPDATE company SET name = '%s' WHERE id = %d", req.Name, req.ID)
	}

	cs.mu.Lock()
	rows, err := cs.db.Query(fmt.Sprintf("SELECT * FROM company WHERE id = %d", req.ID))
	if err != nil {
		cs.l.Printf("%s occured while executing UpdateCompany SQL query", err.Error())
		return nil, err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		e := &accounting.CompanyResponce{}
		err := rows.Scan(&e.ID, &e.Name, &e.LegalForm)
		if err != nil {
			cs.l.Println(err)
			return nil, err
		}
		found = true
	}
	if !found {
		return nil, errors.New("Company not found")
	}
	_, _ = cs.db.Exec(query)
	cs.mu.Unlock()

	return &accounting.CompanyResponce{StatusCode: 200}, nil
}

// GetCompany gets company by id
func (cs *CompanyServer) GetCompany(ctx context.Context, req *accounting.CompanyRequest) (*accounting.CompanyResponce, error) {
	query := fmt.Sprintf("SELECT * FROM company WHERE id = %d", req.ID)

	cs.mu.Lock()
	rows, err := cs.db.Query(query)
	cs.mu.Unlock()
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

	return nil, errors.New("Company not found")
}

// GetCompanyEmployees gets company employees for specified company id
func (cs *CompanyServer) GetCompanyEmployees(ctx context.Context, req *accounting.CompanyRequest) (*accounting.CompanyEmployeesResponce, error) {
	query := fmt.Sprintf("SELECT * FROM employee WHERE company_id = %d", req.ID)

	cs.mu.Lock()
	rows, err := cs.db.Query(query)
	cs.mu.Unlock()
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

// DeleteCompany deletes company by id
func (cs *CompanyServer) DeleteCompany(ctx context.Context, req *accounting.CompanyRequest) (*accounting.CompanyResponce, error) {
	query := fmt.Sprintf("DELETE FROM company WHERE id = %d", req.ID)

	cs.mu.Lock()
	res, err := cs.db.Exec(query)
	cs.mu.Unlock()
	if err != nil {
		cs.l.Printf("%s occured while executing DeleteCompany SQL query", err.Error())
		return nil, err
	}

	if n, err := res.RowsAffected(); err != nil || n == 0 {
		return nil, errors.New("Company not found")
	}
	return &accounting.CompanyResponce{StatusCode: 200}, nil
}
