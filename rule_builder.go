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
	Message   string `json:"message"`
	FieldName string `json:"field"`
}

// Error message
func (v validationError) Error() string {
	return v.Message
}

// Field name
func (v validationError) Field() string {
	return v.FieldName
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

// -----

// Validator to implement a rule
type Validator interface {
	SetName(string)
	Name() string
	HtmlCompatible() bool
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
			err = validator.Validate(subject)
		} else {
			err = validator.Validate(vmap[validator.Name()])
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
		if !v.HtmlCompatible() {
			continue
		}
		name := v.Name()
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
	ft := findStructField(value, fv)
	if ft == nil {
		panic(errors.New("can't find field"))
	}

	tag := strings.Split(ft.Tag.Get("json"), ",")[0]
	if tag == "" {
		tag = ft.Name
	}
	return tag
}

// findStructField looks for a field in the given struct.
// The field being looked for should be a pointer to the actual struct field.
// If found, the field info will be returned. Otherwise, nil will be returned.
func findStructField(structValue reflect.Value, fieldValue reflect.Value) *reflect.StructField {
	ptr := fieldValue.Pointer()
	for i := structValue.NumField() - 1; i >= 0; i-- {
		sf := structValue.Type().Field(i)
		if ptr == structValue.Field(i).UnsafeAddr() {
			// do additional type comparison because it's possible that the address of
			// an embedded struct is the same as the first field of the embedded struct
			if sf.Type == fieldValue.Elem().Type() {
				return &sf
			}
		}
	}
	return nil
}

// joinSentences converts a list of strings to a paragraph
func joinSentences(list []string) string {
	l := len(list)
	if l == 0 {
		return ""
	}
	final := list[0]
	for i := 1; i < l; i++ {
		if i == l-1 {
			final = final + list[i] + "."
		} else {
			final = final + list[i] + ". "
		}
	}
	return final
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
			vmap[name] = f.Interface()
		}
	}
	return vmap
}
