package xvalid

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

// Error when a rule is broken
type Error interface {
	Error() string
	Field() []string
}

// validationError implements Error interface
type validationError struct {
	message string
	field   []string
}

// Error message
func (v validationError) Error() string {
	return v.message
}

// Field name
func (v validationError) Field() []string {
	return v.field
}

func (e validationError) MarshalJSON() ([]byte, error) {
	// only use the last field name for embeded structs
	return json.MarshalIndent(struct {
		Message   string `json:"message"`
		FieldName string `json:"field"`
	}{e.message, jsonFieldName(e.field)}, "", "	")
}

// NewError creates new validation error
func NewError(message string, field []string) Error {
	return &validationError{
		field:   field,
		message: message,
	}
}

// ErrorSlice is a list of Error
type ErrorSlice []Error

// Error will combine all errors into a list of sentences
func (e ErrorSlice) Error() string {
	list := make([]string, len(e))
	for i := range e {
		list[i] = e[i].Error()
	}
	return joinSentences(list)
}

// Unwrap errors
func (e ErrorSlice) Unwrap() []error {
	errs := make([]error, len(e))
	for i := range e {
		errs[i] = e[i]
	}
	return errs
}

// ToMap converts to map
func (e ErrorSlice) ToMap() ErrorMap {
	errs := make(ErrorMap)
	for i, err := range e {
		errs[jsonFieldName(err.Field())] = e[i]
	}
	return errs
}

// ErrorMap is a map of Error
type ErrorMap map[string]Error

// Error will combine all errors into a list of sentences
func (e ErrorMap) Error() string {
	list := make([]string, 0)
	for _, err := range e {
		list = append(list, err.Error())
	}
	return joinSentences(list)
}

// Unwrap errors
func (e ErrorMap) Unwrap() []error {
	errs := make([]error, 0)
	for _, err := range e {
		errs = append(errs, err)
	}
	return errs
}

// ToSlice converts to slice
func (e ErrorMap) ToSlice() ErrorSlice {
	errs := make(ErrorSlice, 0)
	for _, err := range e {
		errs = append(errs, err)
	}
	return errs
}

// -----

// Validator to implement a rule
type Validator interface {
	SetField(...string)
	Field() []string
	CanExport() bool
	SetMessage(string) Validator
	Validate(any) Error
}

// Rules for creating a chain of rules for validating a struct
type Rules struct {
	validators []Validator
	structPtr  any
}

// New rule chain
func New(structPtr any) Rules {
	return Rules{
		structPtr:  structPtr,
		validators: make([]Validator, 0),
	}
}

// Field adds validators for a field
func (r Rules) Field(fieldPtr any, validators ...Validator) Rules {
	for _, validator := range validators {
		validator.SetField(getField(r.structPtr, fieldPtr)...)
		r.validators = append(r.validators, validator)
	}
	return r
}

// Struct adds validators for the struct
func (r Rules) Struct(validators ...Validator) Rules {
	r.validators = append(r.validators, validators...)
	return r
}

// Validate a struct and return Errors
func (r Rules) Validate(subject any) ErrorSlice {
	errs := make(ErrorSlice, 0)
	vmap := structToMap(subject)
	for _, validator := range r.validators {
		var err Error
		if validator.Field() == nil || len(validator.Field()) == 0 {
			// struct validation
			err = validator.Validate(subject)
		} else {
			// field validation
			v := vmap
			for _, p := range validator.Field() {
				switch v2 := v[p].(type) {
				default:
					err = validator.Validate(v2)
				case map[string]any:
					v = v2
				}
			}
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Validators for this chain
func (r Rules) Validators() []Validator {
	return r.validators
}

func (r Rules) MarshalJSON() ([]byte, error) {
	rmap := make(map[string][]any)
	validators := r.Validators()
	for _, v := range validators {
		if !v.CanExport() {
			continue
		}
		name := jsonFieldName(v.Field())
		rules, ok := rmap[name]
		if !ok {
			rules = make([]any, 0)
		}
		rules = append(rules, v)
		rmap[name] = rules
	}
	return json.MarshalIndent(rmap, "", "	")
}

// -------------------

func getField(structPtr any, fieldPtr any) []string {
	value := reflect.ValueOf(structPtr)
	if value.Kind() != reflect.Ptr || !value.IsNil() && value.Elem().Kind() != reflect.Struct {
		panic(errors.New("struct is not pointer"))
	}
	if value.IsNil() {
		panic(errors.New("struct is nil"))
	}
	value = value.Elem()

	fv := reflect.ValueOf(fieldPtr)
	if fv.Kind() != reflect.Ptr {
		panic(errors.New("field is not pointer"))
	}
	fields := findStructField(value, fv, make([]*reflect.StructField, 0))
	if len(fields) == 0 {
		panic(errors.New("can't find field"))
	}

	parts := make([]string, 0)
	for _, f := range fields {
		tag := strings.Split(f.Tag.Get("json"), ",")[0]
		if tag == "" {
			tag = f.Name
		}
		parts = append(parts, tag)
	}
	return parts
}

// findStructField looks for a field in the given struct.
// The field being looked for should be a pointer to the actual struct field.
// If found, the fields will be returned. Otherwise, an empty list will be returned.
func findStructField(structValue reflect.Value, fieldValue reflect.Value, results []*reflect.StructField) []*reflect.StructField {
	ptr := fieldValue.Pointer()
	depth := len(results)
	for i := structValue.NumField() - 1; i >= 0; i-- {
		sf := structValue.Type().Field(i)
		if ptr == structValue.Field(i).UnsafeAddr() {
			if sf.Anonymous {
				return findStructField(structValue.Field(i), fieldValue, append(results, &sf))
			}
			return append(results, &sf)
		} else if sf.Anonymous {
			tmp := findStructField(structValue.Field(i), fieldValue, append(results, &sf))
			if len(tmp) > depth+1 {
				return tmp
			}
		}
	}
	return results
}

// joinSentences converts a list of strings to a paragraph
func joinSentences(list []string) string {
	l := len(list)
	if l == 0 {
		return ""
	}
	return strings.Join(list, ". ") + "."
}

// structToMap converts struct to map and uses the json name if available
func structToMap(structPtr any) map[string]any {
	vmap := make(map[string]any)
	structValue := reflect.ValueOf(structPtr)
	for i := structValue.NumField() - 1; i >= 0; i-- {
		sf := structValue.Type().Field(i)
		name := strings.Split(sf.Tag.Get("json"), ",")[0]
		if name == "" {
			name = sf.Name
		}
		f := structValue.Field(i)
		if f.CanInterface() {
			if sf.Anonymous {
				vmap[name] = structToMap(f.Interface())
			} else {
				vmap[name] = f.Interface()
			}
		}
	}
	return vmap
}

// jsonFieldName returns the last field name
func jsonFieldName(field []string) string {
	if field == nil {
		return ""
	}
	return field[len(field)-1]
}
