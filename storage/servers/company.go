package servers

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Felley/accounting-service/protos/accounting"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CompanyServer ...
type CompanyServer struct {
	db *sql.DB
	accounting.UnimplementedCompanyAccountingServer
}

// NewCompanyServer returns new employee storage processing server
func NewCompanyServer(db *sql.DB) *CompanyServer {
	return &CompanyServer{db: db}
}

// AddCompany ...
func (es *CompanyServer) AddCompany(context.Context, *accounting.CompanyRequest) (*accounting.CompanyResponce, error) {
	fmt.Println("Hi from AddCompany!")
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCompany not implemented")
}

// UpdateCompany ...
func (es *CompanyServer) UpdateCompany(context.Context, *accounting.CompanyRequest) (*accounting.CompanyResponce, error) {
	fmt.Println("Hi from UpdateCompany!")
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCompany not implemented")
}

// GetCompany ...
func (es *CompanyServer) GetCompany(context.Context, *accounting.CompanyRequest) (*accounting.CompanyResponce, error) {
	fmt.Println("Hi from GetCompany!")
	return nil, status.Errorf(codes.Unimplemented, "method GetCompany not implemented")
}

// DeleteCompany ...
func (es *CompanyServer) DeleteCompany(context.Context, *accounting.CompanyRequest) (*accounting.CompanyResponce, error) {
	fmt.Println("Hi from DeleteCompany!")
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCompany not implemented")
}
