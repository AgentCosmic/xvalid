package xvalid

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"

	"golang.org/x/exp/constraints"
)

//
// ==================== Required ====================
//

// RequiredValidator field must not be zero
type RequiredValidator struct {
	field   []string
	message string
}

// Field of the field
func (c *RequiredValidator) Field() []string {
	return c.field
}

// SetField of the field
func (c *RequiredValidator) SetField(field ...string) {
	c.field = field
}

// SetMessage set error message
func (c *RequiredValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// Validate the value
func (c *RequiredValidator) Validate(value any) Error {
	v := reflect.ValueOf(value)
	zero := false
	kind := v.Kind()
	if !v.IsValid() {
		zero = true
	} else if v.IsZero() {
		zero = true
	} else if (kind == reflect.Ptr || kind == reflect.Interface) && v.Elem().IsZero() {
		zero = true
	} else if (kind == reflect.Array || kind == reflect.Slice || kind == reflect.Map) && v.Len() == 0 {
		zero = true
	}
	if zero {
		return createError(c.field, c.message, fmt.Sprintf("Please enter the %v", jsonFieldName(c.field)))
	}
	return nil
}

// MarshalJSON for this validator
func (c *RequiredValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rule    string `json:"rule"`
		Message string `json:"message,omitempty"`
	}{"required", c.message})
}

// CanExport for this validator
func (c *RequiredValidator) CanExport() bool {
	return true
}

// Required fields must not be zero
func Required() *RequiredValidator {
	return &RequiredValidator{}
}

//
// ==================== MinLength ====================
//

// MinLengthValidator field must have minimum length
type MinLengthValidator struct {
	field    []string
	message  string
	min      int64
	optional bool
}

// Field of the field
func (c *MinLengthValidator) Field() []string {
	return c.field
}

// SetField of the field
func (c *MinLengthValidator) SetField(name ...string) {
	c.field = name
}

// SetMessage set error message
func (c *MinLengthValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// SetOptional don't validate if the value is zero
func (c *MinLengthValidator) SetOptional() Validator {
	c.optional = true
	return c
}

// Validate the value
func (c *MinLengthValidator) Validate(value any) Error {
	str, ok := value.(string)
	if !ok {
		if c.optional {
			return nil
		} else {
			return createError(c.field, c.message, fmt.Sprintf("Please lengthen %s to %d characters or more", jsonFieldName(c.field), c.min))
		}
	}
	if c.optional && str == "" {
		return nil
	}
	if len([]rune(str)) < int(c.min) {
		return createError(c.field, c.message, fmt.Sprintf("Please lengthen %s to %d characters or more", jsonFieldName(c.field), c.min))
	}
	return nil
}

// MarshalJSON for this validator
func (c *MinLengthValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rule    string `json:"rule"`
		Min     int64  `json:"min"`
		Message string `json:"message,omitempty"`
	}{"minLength", c.min, c.message})
}

// CanExport for this validator
func (c *MinLengthValidator) CanExport() bool {
	return true
}

// MinLength field must have minimum length
func MinLength(min int64) *MinLengthValidator {
	return &MinLengthValidator{
		min: min,
	}
}

//
// ==================== MaxLength ====================
//

// MaxLengthValidator field have maximum length
type MaxLengthValidator struct {
	ifeld   []string
	message string
	max     int64
}

// Field of the field
func (c *MaxLengthValidator) Field() []string {
	return c.ifeld
}

// SetField of the field
func (c *MaxLengthValidator) SetField(name ...string) {
	c.ifeld = name
}

// SetMessage set error message
func (c *MaxLengthValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// Validate the value
func (c *MaxLengthValidator) Validate(value any) Error {
	v, ok := value.(string)
	if !ok {
		return nil
	}
	if len([]rune(v)) > int(c.max) {
		return createError(c.ifeld, c.message, fmt.Sprintf("Please shorten %s to %d characters or less", jsonFieldName(c.ifeld), c.max))
	}
	return nil
}

// MarshalJSON for this validator
func (c *MaxLengthValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rule    string `json:"rule"`
		Max     int64  `json:"max"`
		Message string `json:"message,omitempty"`
	}{"maxLength", c.max, c.message})
}

// CanExport for this validator
func (c *MaxLengthValidator) CanExport() bool {
	return true
}

// MaxLength field have maximum length
func MaxLength(max int64) *MaxLengthValidator {
	return &MaxLengthValidator{
		max: max,
	}
}

//
// ==================== Min ====================
//

// MinValidator field have minimum value
type MinValidator struct {
	field    []string
	message  string
	min      int64
	optional bool
}

// Field of the field
func (c *MinValidator) Field() []string {
	return c.field
}

// SetField of the field
func (c *MinValidator) SetField(name ...string) {
	c.field = name
}

// SetMessage set error message
func (c *MinValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// SetOptional don't validate if the value is zero
func (c *MinValidator) SetOptional() Validator {
	c.optional = true
	return c
}

// Validate the value
func (c *MinValidator) Validate(value any) Error {
	rv := reflect.ValueOf(value)
	newError := func() Error {
		return createError(c.field, c.message, fmt.Sprintf("Please increase %s to be %v or more", jsonFieldName(c.field), c.min))
	}
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if isLess(toInt64(value), c.min, c.optional) {
			return newError()
		}
	case reflect.Float32, reflect.Float64:
		if isLess(toFloat64(value), float64(c.min), c.optional) {
			return newError()
		}
	case reflect.Invalid:
		if !c.optional {
			return newError()
		}
	default:
		panic(fmt.Errorf("type not supported: %v", rv.Type()))
	}
	return nil
}

// MarshalJSON for this validator
func (c *MinValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rule    string `json:"rule"`
		Min     int64  `json:"min"`
		Message string `json:"message,omitempty"`
	}{"min", c.min, c.message})
}

// CanExport for this validator
func (c *MinValidator) CanExport() bool {
	return true
}

// Min field have minimum value
func Min(min int64) *MinValidator {
	return &MinValidator{
		min: min,
	}
}

//
// ==================== Max ====================
//

// MaxValidator field have maximum value
type MaxValidator struct {
	field   []string
	message string
	max     int64
}

// Field of the field
func (c *MaxValidator) Field() []string {
	return c.field
}

// SetField of the field
func (c *MaxValidator) SetField(name ...string) {
	c.field = name
}

// SetMessage set error message
func (c *MaxValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// Validate the value
func (c *MaxValidator) Validate(value any) Error {
	rv := reflect.ValueOf(value)
	newError := func() Error {
		return createError(c.field, c.message, fmt.Sprintf("Please decrease %s to be %v or less", jsonFieldName(c.field), c.max))
	}
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if isMore(toInt64(value), c.max) {
			return newError()
		}
	case reflect.Float32, reflect.Float64:
		if isMore(toFloat64(value), float64(c.max)) {
			return newError()
		}
	case reflect.Invalid:
		return nil
	default:
		panic(fmt.Errorf("type not supported: %v", rv.Type()))
	}
	return nil
}

// MarshalJSON for this validator
func (c *MaxValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rule    string `json:"rule"`
		Max     int64  `json:"max"`
		Message string `json:"message,omitempty"`
	}{"max", c.max, c.message})
}

// CanExport for this validator
func (c *MaxValidator) CanExport() bool {
	return true
}

// Max field have maximum value
func Max(max int64) *MaxValidator {
	return &MaxValidator{
		max: max,
	}
}

//
// ==================== Pattern ====================
//

// PatternValidator field must match regexp
type PatternValidator struct {
	field    []string
	message  string
	re       *regexp.Regexp
	optional bool
}

// Field of the field
func (c *PatternValidator) Field() []string {
	return c.field
}

// SetField of the field
func (c *PatternValidator) SetField(name ...string) {
	c.field = name
}

// SetMessage set error message
func (c *PatternValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// SetOptional don't validate if the value is zero
func (c *PatternValidator) SetOptional() Validator {
	c.optional = true
	return c
}

// Validate the value
func (c *PatternValidator) Validate(value any) Error {
	str, ok := value.(string)
	if !ok {
		if c.optional {
			return nil
		} else {
			return createError(c.field, c.message, fmt.Sprintf("Please correct %s into a valid format", jsonFieldName(c.field)))
		}
	}
	if c.optional && str == "" {
		return nil
	}
	if c.re.MatchString(str) {
		return nil
	}
	return createError(c.field, c.message, fmt.Sprintf("Please correct %s into a valid format", jsonFieldName(c.field)))
}

// MarshalJSON for this validator
func (c *PatternValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rule    string `json:"rule"`
		Pattern string `json:"pattern"`
		Message string `json:"message,omitempty"`
	}{"pattern", c.re.String(), c.message})
}

// CanExport for this validator
func (c *PatternValidator) CanExport() bool {
	return true
}

// Pattern field must match regexp
func Pattern(pattern string) *PatternValidator {
	return &PatternValidator{
		re: regexp.MustCompile(pattern),
	}
}

//
// ==================== Email ====================
//

// EmailValidator field must be a valid email address
type EmailValidator struct {
	Validator
	field    []string
	message  string
	optional bool
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Email field must be a valid email address
func Email() *EmailValidator {
	return &EmailValidator{}
}

// Field of the field
func (c *EmailValidator) Field() []string {
	return c.field
}

// SetField of the field
func (c *EmailValidator) SetField(name ...string) {
	c.field = name
}

// SetMessage set error message
func (c *EmailValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// SetOptional don't validate if the value is zero
func (c *EmailValidator) SetOptional() Validator {
	c.optional = true
	return c
}

// Validate the value
func (c *EmailValidator) Validate(value any) Error {
	str, ok := value.(string)
	if !ok {
		if c.optional {
			return nil
		} else {
			return createError(c.field, c.message, fmt.Sprintf("Please use a valid email address for %s", jsonFieldName(c.field)))
		}
	}
	if c.optional && str == "" {
		return nil
	}
	if emailRegex.MatchString(str) {
		return nil
	}
	return createError(c.field, c.message, fmt.Sprintf("Please use a valid email address for %s", jsonFieldName(c.field)))
}

// CanExport for this validator
func (c *EmailValidator) CanExport() bool {
	return true
}

// MarshalJSON for this validator
func (c *EmailValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rule    string `json:"rule"`
		Type    string `json:"type"`
		Pattern string `json:"pattern"`
		Message string `json:"message,omitempty"`
	}{"type", "email", emailRegex.String(), c.message})
}

// IsEmail returns true if the string is an email
func IsEmail(email string) bool {
	return emailRegex.MatchString(email)
}

//
// ==================== Options ====================
//

// OptionsValidator for whitelisting accepted values
type OptionsValidator struct {
	field   []string
	message string
	options []any
}

// Field of the field
func (c *OptionsValidator) Field() []string {
	return c.field
}

// SetField of the field
func (c *OptionsValidator) SetField(name ...string) {
	c.field = name
}

// SetMessage set error message
func (c *OptionsValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// Validate the value
func (c *OptionsValidator) Validate(value any) Error {
	v := reflect.ValueOf(value)
	actual := v.Interface()
	for _, opt := range c.options {
		if opt == actual {
			return nil
		}
	}
	return createError(c.field, c.message, fmt.Sprintf("Please select one of the valid options for %s", jsonFieldName(c.field)))
}

// CanExport for this validator
func (c *OptionsValidator) CanExport() bool {
	return true
}

// MarshalJSON for this validator
func (c *OptionsValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rule    string `json:"rule"`
		Options []any  `json:"options"`
		Message string `json:"message,omitempty"`
	}{"options", c.options, c.message})
}

// Options for whitelisting accepted values
func Options(options ...any) Validator {
	return &OptionsValidator{
		options: options,
	}
}

//
// ==================== FieldFunc ====================
//

// FieldFuncValidator for validating with custom function
type FieldFuncValidator struct {
	field   []string
	message string
	checker func([]string, any) Error
}

// Field of the field
func (c *FieldFuncValidator) Field() []string {
	return c.field
}

// SetField of the field
func (c *FieldFuncValidator) SetField(name ...string) {
	c.field = name
}

// SetMessage set error message
func (c *FieldFuncValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// Validate the value
func (c *FieldFuncValidator) Validate(value any) Error {
	return c.checker(c.field, value)
}

// CanExport for this validator
func (c *FieldFuncValidator) CanExport() bool {
	return false
}

// FieldFunc for validating with custom function
func FieldFunc(f func([]string, any) Error) Validator {
	return &FieldFuncValidator{
		checker: f,
	}
}

//
// ==================== StructFunc ====================
//

// StructFuncValidator validate struct with custom function. Add to rules with .Struct().
type StructFuncValidator struct {
	field   []string
	message string
	checker func(any) Error
}

// Field of the field
func (c *StructFuncValidator) Field() []string {
	return c.field
}

// SetField of the field
func (c *StructFuncValidator) SetField(name ...string) {
	c.field = name
}

// SetMessage set error message
func (c *StructFuncValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// Validate the value
func (c *StructFuncValidator) Validate(value any) Error {
	return c.checker(value)
}

// CanExport for this validator
func (c *StructFuncValidator) CanExport() bool {
	return false
}

// StructFunc validate struct with custom function
func StructFunc(f func(any) Error) Validator {
	return &StructFuncValidator{
		checker: f,
	}
}

//
// ====================
//

func createError(field []string, custom string, fallback string) Error {
	if custom != "" {
		return NewError(custom, field)
	}
	return NewError(fallback, field)
}

func toInt64(value any) int64 {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	}
	panic(fmt.Errorf("cannot convert %v to int64", v.Kind()))
}

func toFloat64(value any) float64 {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return v.Float()
	}
	panic(fmt.Errorf("cannot convert %v to float64", v.Kind()))
}

func isLess[T number](value T, min T, optional bool) bool {
	if optional && value == 0 {
		return false
	}
	if value < min {
		return true
	}
	return false
}

func isMore[T number](value T, max T) bool {
	return value > max
}

type number interface {
	constraints.Integer | constraints.Float
}
