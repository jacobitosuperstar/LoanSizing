// Error structs for the package

package financial_formulas

import "fmt"


type ValidationError struct {
    Field string
    Value any
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("Validation Error\nField: %v\nValue: %v\n%v", e.Field, e.Value, e.Message)
}

type ValueError struct {
    Field string
    Value any
    Message string
}

func (e *ValueError) Error() string {
    return fmt.Sprintf("Value Error\nField: %v\nValue: %v\n%v", e.Field, e.Value, e.Message)
}
