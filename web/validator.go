package web

// Based on:
// Let's Go Further - Chapter 4.5 (Alex Edwards)
type Validator struct {
	Errors map[string]string
}

func NewValidator() *Validator {
	v := Validator{
		Errors: make(map[string]string),
	}
	return &v
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}
