package data

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

// Company defines the structure for API
type Company struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	LegalForm string `json:"legalForm" validate:"legalForm"`
}

// FromJSON unmarshalls []bytes to Company struct
func (e *Company) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(e)
}

// Validate function validates incoming legalForm field from JSON
func (e *Company) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("legalForm", validateLegalForm, false)
	return validate.Struct(e)
}

func validateLegalForm(fl validator.FieldLevel) bool {
	return true
}
