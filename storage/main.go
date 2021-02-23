package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"

	"github.com/Felley/accounting-service/protos/accounting"
	"github.com/Felley/accounting-service/storage/servers"
	"github.com/Felley/accounting-service/storage/tables"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

// Starting storage processing service
func main() {
	gc := grpc.NewServer()

	db, err := sql.Open("mysql", "root:root@tcp(0.0.0.0:3306)/accounting_db")
	if err != nil {
		panic(err)
	}
	_ = tables.CreateEmployeeTable(db)
	_ = tables.CreateCompanyTable(db)
	_ = tables.ListTables(db)
	es := servers.NewEmployeeServer(db)
	cs := servers.NewCompanyServer(db)
	defer db.Close()

	accounting.RegisterEmployeeAccountingServer(gc, es)
	accounting.RegisterCompanyAccountingServer(gc, cs)

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		fmt.Printf("Error: %e", err)
		os.Exit(1)
	}
	gc.Serve(l)
}
