package data

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

// Company defines the structure for API
type Company struct {
	ID        int64  `json:"id" validate:"omitempty,numeric"`
	Name      string `json:"name" validate:"required"`
	LegalForm string `json:"legalForm" validate:"omitempty,legalForm"`
}

// NewCompany creates new Employee struct from incoming data
func NewCompany(id int64, name string, legalForm string) *Company {
	return &Company{
		ID:        id,
		Name:      name,
		LegalForm: legalForm,
	}
}

// FromJSON unmarshalls []bytes to Company struct
func (e *Company) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(e)
}

// Validate function validates incoming JSON fields
func (e *Company) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("legalForm", validateLegalForm, false)
	return validate.Struct(e)
}

func validateLegalForm(fl validator.FieldLevel) bool {
	return true
}
