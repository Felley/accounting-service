package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Felley/accounting-service/api/handlers"
	"github.com/Felley/accounting-service/protos/accounting"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

// Starting point for API service
func main() {
	logger := log.New(os.Stdout, "accounting-api", log.LstdFlags)
	employeeLogger := log.New(os.Stdout, "employee-api", log.LstdFlags)
	companyLogger := log.New(os.Stdout, "company-api", log.LstdFlags)

	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	router := mux.NewRouter()

	// create gRPC client
	employeeClient := accounting.NewEmployeeAccountingClient(conn)
	companyClient := accounting.NewCompanyAccountingClient(conn)

	// create the handlers
	employeeHandler := handlers.NewEmployeeHandler(employeeLogger, employeeClient)
	companyHandler := handlers.NewCompanyHandler(companyLogger, companyClient)

	// configure employee routes
	postEmployeeRouter := router.Methods(http.MethodPost).Subrouter()
	postEmployeeRouter.HandleFunc("/employee", employeeHandler.AddEmployee)
	postEmployeeRouter.Use(employeeHandler.MiddlewareEmployeeValidation)

	putEmployeeRouter := router.Methods(http.MethodPut).Subrouter()
	putEmployeeRouter.HandleFunc("/employee", employeeHandler.UpdateEmployee)
	putEmployeeRouter.Use(employeeHandler.MiddlewareEmployeeValidation)

	getEmployeeRouter := router.Methods(http.MethodGet).Subrouter()
	getEmployeeRouter.HandleFunc("/employee/{id:[0-9]+}", employeeHandler.GetEmployee)

	postEmployeeFormRouter := router.Methods(http.MethodPost).Subrouter()
	postEmployeeFormRouter.HandleFunc("/employee/{id:[0-9]+}", employeeHandler.PostFormEmployee)

	deleteEmployeeRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteEmployeeRouter.HandleFunc("/employee/{id:[0-9]+}", employeeHandler.DeleteEmployee)

	// configure company routes
	postCompanyRouter := router.Methods(http.MethodPost).Subrouter()
	postCompanyRouter.HandleFunc("/company/", companyHandler.AddCompany)
	postCompanyRouter.Use(companyHandler.MiddlewareCompanyValidation)

	putCompanyRouter := router.Methods(http.MethodPut).Subrouter()
	putCompanyRouter.HandleFunc("/company/", companyHandler.UpdateCompany)
	putCompanyRouter.Use(companyHandler.MiddlewareCompanyValidation)

	getCompanyRouter := router.Methods(http.MethodGet).Subrouter()
	getCompanyRouter.HandleFunc("/company/{id:[0-9]+}", companyHandler.GetCompany)

	postCompanyFormRouter := router.Methods(http.MethodPost).Subrouter()
	postCompanyFormRouter.HandleFunc("/company/{id:[0-9]+}", companyHandler.PostFormCompany)

	deleteCompanyRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteCompanyRouter.HandleFunc("/company/{id:[0-9]+}", companyHandler.DeleteCompany)

	getCompanyEmpoyeesRouter := router.Methods(http.MethodGet).Subrouter()
	getCompanyEmpoyeesRouter.HandleFunc("/company/{id:[0-9]+}/employees", companyHandler.GetCompanyEmployees)

	server := http.Server{
		Addr:    ":9090",
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	interruptionChan := make(chan os.Signal)
	signal.Notify(interruptionChan, os.Interrupt)
	signal.Notify(interruptionChan, os.Kill)

	sig := <-interruptionChan
	logger.Println("Recieved server terminate, shutdown", sig)

	timeout, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeout)
}
