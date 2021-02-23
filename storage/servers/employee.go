package servers

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Felley/accounting-service/protos/accounting"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EmployeeServer ...
type EmployeeServer struct {
	db *sql.DB
	accounting.UnimplementedEmployeeAccountingServer
}

// NewEmployeeServer returns new employee storage processing server
func NewEmployeeServer(db *sql.DB) *EmployeeServer {
	return &EmployeeServer{db: db}
}

// AddEmployee ...
func (es *EmployeeServer) AddEmployee(context.Context, *accounting.EmployeeRequest) (*accounting.EmployeeResponce, error) {
	fmt.Println("Hi from AddEmployee!")
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEmployee not implemented")
}

// UpdateEmployee ...
func (es *EmployeeServer) UpdateEmployee(context.Context, *accounting.EmployeeRequest) (*accounting.EmployeeResponce, error) {
	fmt.Println("Hi from UpdateEmployee!")
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEmployee not implemented")
}

// GetEmployee ...
func (es *EmployeeServer) GetEmployee(context.Context, *accounting.EmployeeRequest) (*accounting.EmployeeResponce, error) {
	fmt.Println("Hi from GetEmployee!")
	return nil, status.Errorf(codes.Unimplemented, "method GetEmployee not implemented")
}

// DeleteEmployee ...
func (es *EmployeeServer) DeleteEmployee(context.Context, *accounting.EmployeeRequest) (*accounting.EmployeeResponce, error) {
	fmt.Println("Hi from DeleteEmployee!")
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEmployee not implemented")
}
