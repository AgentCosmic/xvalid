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
	Field() string
}

// validationError implements Error interface
type validationError struct {
	Message   string
	FieldName string
}

// Error message
func (v validationError) Error() string {
	return v.Message
}

// Field name
func (v validationError) Field() string {
	return v.FieldName
}

func (e validationError) MarshalJSON() ([]byte, error) {
	// only use the last field name for embeded structs
	return json.MarshalIndent(struct {
		Message   string `json:"message"`
		FieldName string `json:"field"`
	}{e.Message, jsonFieldName(e.FieldName)}, "", "	")
}

// NewError creates new validation error
func NewError(message, fieldName string) Error {
	return &validationError{
		FieldName: fieldName,
		Message:   message,
	}
}

// Errors is a list of Error
type Errors []Error

// Error will combine all errors into a list of sentences
func (v Errors) Error() string {
	list := make([]string, len(v))
	for i := range v {
		list[i] = v[i].Error()
	}
	return joinSentences(list)
}

// Unwrap errors
func (v Errors) Unwrap() []error {
	errs := make([]error, len(v))
	for _, e := range v {
		errs = append(errs, e)
	}
	return errs
}

// -----

// Validator to implement a rule
type Validator interface {
	SetName(string)
	Name() string
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
		validator.SetName(getFieldName(r.structPtr, fieldPtr))
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
func (r Rules) Validate(subject any) error {
	errs := make(Errors, 0)
	vmap := structToMap(subject)
	for _, validator := range r.validators {
		var err Error
		if validator.Name() == "" {
			// struct validation
			err = validator.Validate(subject)
		} else {
			// field validation
			path := strings.Split(validator.Name(), ".")
			v := vmap
			for _, p := range path {
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
		name := jsonFieldName(v.Name())
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

func getFieldName(structPtr any, fieldPtr any) string {
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
	return strings.Join(parts, ".")
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

func jsonFieldName(field string) string {
	parts := strings.Split(field, ".")
	return parts[len(parts)-1]
}
