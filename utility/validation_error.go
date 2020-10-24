package utility

import "fmt"

type ValidationError struct {
	ErrorMessages []error `json:"errorMessages"`
}

func (ve ValidationError) Error() string {
	totalErrors := len(ve.ErrorMessages)
	if totalErrors == 0 {
		return "No validation error found"
	}
	if totalErrors == 1 {
		return fmt.Sprintf("%v validation error found", totalErrors)
	}
	return fmt.Sprintf("%v validation errors found", totalErrors)
}
