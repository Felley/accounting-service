package data

import (
	"errors"
	"testing"
)

func TestCheckCompanyValidation(t *testing.T) {
	t.Parallel()
	var tableDrivenTests = []struct {
		testName string
		in       *Company
		expected error
	}{
		{
			"Empty",
			&Company{},
			errors.New("Empty struct"),
		},
		{
			"Correct",
			&Company{
				ID:        1,
				Name:      "Test",
				LegalForm: "OAO",
			},
			nil,
		},
		{
			"No name",
			&Company{
				ID:        1,
				LegalForm: "OAO",
			},
			errors.New("Name is not defined"),
		},
		{
			"Correct with only name",
			&Company{
				Name: "Test",
			},
			nil,
		},
	}

	for _, test := range tableDrivenTests {
		t.Run(test.testName, func(t *testing.T) {
			err := test.in.Validate()

			if test.expected == nil && err == nil {

			} else {
				if test.expected != nil && err != nil {

				} else {
					t.Errorf("got %q, want %q", err, test.expected)
				}
			}
		})
	}
}
