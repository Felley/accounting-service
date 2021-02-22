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

	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	router := mux.NewRouter()

	// create gRPC client
	employeeClient := accounting.NewEmployeeAccountingClient(conn)

	// create the handlers
	handler := handlers.NewEmployeeHandler(employeeLogger, employeeClient)

	// configure routes
	postEmployeeRouter := router.Methods(http.MethodPost).Subrouter()
	postEmployeeRouter.HandleFunc("/employee", handler.AddEmployee)
	postEmployeeRouter.Use(handler.MiddlewareEmployeeValidation)

	putEmployeeRouter := router.Methods(http.MethodPut).Subrouter()
	putEmployeeRouter.HandleFunc("/employee/{id:[0-9]+}", handler.UpdateEmployee)
	putEmployeeRouter.Use(handler.MiddlewareEmployeeValidation)

	getEmployeeRouter := router.Methods(http.MethodGet).Subrouter()
	getEmployeeRouter.HandleFunc("/employee/{id:[0-9]+}", handler.GetEmployee)

	postCompanyEmployeeRouter := router.Methods(http.MethodPost).Subrouter()
	postCompanyEmployeeRouter.HandleFunc("/employee/{id:[0-9]+}", handler.PostFormEmployee)

	deleteEmployeeRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteEmployeeRouter.HandleFunc("/employee/{id:[0-9]+}", handler.DeleteEmployee)

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
