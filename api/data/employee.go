package data

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

// Employee defines the structure for API
type Employee struct {
	ID         int64  `json:"id" validate:"omitempty,numeric"`
	Name       string `json:"name" validate:"required,alphaunicode"`
	SecondName string `json:"secondName" validate:"omitempty,alphaunicode"`
	Surname    string `json:"surname" validate:"omitempty,alphaunicode"`
	HireDate   string `json:"hireDate" validate:"omitempty,datetime=2006-01-02"`
	Position   string `json:"position" validate:"omitempty,position"`
	CompanyID  int64  `json:"companyId" validate:"omitempty,numeric"`
}

// NewEmployee creates new Employee struct from incoming data
func NewEmployee(id int64, name string, secondName string, surname string, hireDate string, position string, companyID int64) *Employee {
	return &Employee{
		ID:         id,
		Name:       name,
		SecondName: secondName,
		Surname:    surname,
		HireDate:   hireDate,
		Position:   position,
		CompanyID:  companyID,
	}
}

// FromJSON unmarshalls []bytes to Employee struct
func (e *Employee) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(e)
}

// Validate function validates incoming JSON fields
func (e *Employee) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("position", validatePosition, false)
	return validate.Struct(e)
}

func validatePosition(fl validator.FieldLevel) bool {
	switch fl.Field().String() {
	case "developer", "manager", "quality assurance", "business analyst", "":
		return true
	default:
		return false
	}
}
