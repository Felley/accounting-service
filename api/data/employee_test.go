package data

import (
	"errors"
	"testing"
)

func TestCheckValidation(t *testing.T) {
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
