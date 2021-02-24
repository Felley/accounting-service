package tables

import (
	"database/sql"
	"fmt"
)

// CreateEmployeeTable creates table employee if it does not exist with id, name, second_name, surname, hire_date, position, company_id fields
func CreateEmployeeTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS employee(id int primary key auto_increment, name text NOT NULL, second_name text, surname text, hire_date date, position text, company_id int)`
	stmt, err := db.Prepare(query)
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Table created successfully..")
	}
	return nil
}

// CreateCompanyTable creates table company with id, name, legal_form fields
func CreateCompanyTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS company(id int primary key auto_increment, name text NOT NULL, legal_form text)`
	stmt, err := db.Prepare(query)
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Table created successfully..")
	}
	return nil
}

// ListTables shows all db tables in console
func ListTables(db *sql.DB) error {
	res, _ := db.Query("SHOW TABLES")

	var table string

	for res.Next() {
		res.Scan(&table)
		fmt.Println(table)
	}
	return nil
}
