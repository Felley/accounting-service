package tables

import (
	"database/sql"
	"fmt"
)

// CreateEmployeeTable ...
func CreateEmployeeTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS employee(id int primary key auto_increment, name text, second_name text, surname text, hire_date date, position text, company_id int)`
	stmt, err := db.Prepare(query)
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Table created successfully..")
	}
	return nil
}

// CreateCompanyTable ...
func CreateCompanyTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS company(id int primary key auto_increment, name text, legal_form text)`
	stmt, err := db.Prepare(query)
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Table created successfully..")
	}
	return nil
}

// ListTables ...
func ListTables(db *sql.DB) error {
	res, _ := db.Query("SHOW TABLES")

	var table string

	for res.Next() {
		res.Scan(&table)
		fmt.Println(table)
	}
	return nil
}
