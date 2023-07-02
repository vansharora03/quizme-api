package validator

type Validator struct {
    Errors map[string]string
}

// Valid returns true if there are no errors in the validation
func (v *Validator) Valid() bool {
    return len(v.Errors) == 0
}

// Check records the validation msg if the condition is not true. 
// After running all necessary calls to this function, the v.Valid() 
// function can be called to see if the validation was successful.
func (v *Validator) Check(condition bool, field, msg string) {
    if !condition {
        v.Add(field, msg)
    }
}

// Adds the msg to field in the validator.
func (v *Validator) Add(field, msg string) {
    v.Errors[field] = msg
}

