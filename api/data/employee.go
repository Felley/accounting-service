package data

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

// Employee defines the structure for API
type Employee struct {
	ID         int64  `json:"id"`
	Name       string `json:"name" validate:"required"`
	SecondName string `json:"secondName"`
	Surname    string `json:"surname"`
	HireDate   string `json:"hireDate" validate:"datetime=2006-01-02"`
	Position   string `json:"position" validate:"position"`
	CompanyID  int64  `json:"companyId"`
}

// FromJSON unmarshalls []bytes to Employee struct
func (e *Employee) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(e)
}

// Validate ...
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
