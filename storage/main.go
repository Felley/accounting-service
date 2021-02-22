package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"

	"github.com/Felley/accounting-service/protos/accounting"
	"github.com/Felley/accounting-service/storage/servers"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

// Starting storage processing service
func main() {
	gc := grpc.NewServer()

	db, err := sql.Open("mysql", "root:root@/accounting_db")
	if err != nil {
		panic(err)
	}

	cs := servers.NewEmployeeServer(db)
	defer db.Close()

	accounting.RegisterEmployeeAccountingServer(gc, cs)

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		fmt.Printf("Error: %e", err)
		os.Exit(1)
	}
	gc.Serve(l)
}
