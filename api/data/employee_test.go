package data

import (
	"errors"
	"testing"
)

func TestCheckEmployeeValidation(t *testing.T) {
	t.Parallel()
	var tableDrivenTests = []struct {
		testName string
		in       *Employee
		expected error
	}{
		{
			"Empty",
			&Employee{},
			errors.New("Empty struct"),
		},
		{
			"Wrong date",
			&Employee{
				Name:     "Aaron",
				HireDate: "2123e",
			},
			errors.New("Invalid date"),
		},
		{
			"No name specified",
			&Employee{
				Surname:  "Aaron",
				HireDate: "2020-09-01",
			},
			errors.New("No name specified"),
		},
		{
			"Correct with only name",
			&Employee{
				Name: "Aaron",
			},
			nil,
		},
		{
			"Invalid Name",
			&Employee{
				Name:     "2",
				HireDate: "2020-09-01",
			},
			errors.New("Name has digits"),
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
