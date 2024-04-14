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
	name    string
	message string
}

// Name of the field
func (c *RequiredValidator) Name() string {
	return c.name
}

// SetName of the field
func (c *RequiredValidator) SetName(name string) {
	c.name = name
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
	} else if (kind == reflect.Array || kind == reflect.Slice) && v.Len() == 0 {
		zero = true
	}
	if zero {
		return createError(c.name, c.message, fmt.Sprintf("Please enter the %v", c.name))
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

// HtmlCompatible for this validator
func (c *RequiredValidator) HtmlCompatible() bool {
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
	name     string
	message  string
	min      int64
	optional bool
}

// Name of the field
func (c *MinLengthValidator) Name() string {
	return c.name
}

// SetName of the field
func (c *MinLengthValidator) SetName(name string) {
	c.name = name
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
	str := value.(string)
	if c.optional && str == "" {
		return nil
	}
	if len([]rune(str)) < int(c.min) {
		return createError(c.name, c.message, fmt.Sprintf("Please lengthen %s to %d characters or more", c.name, c.min))
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

// HtmlCompatible for this validator
func (c *MinLengthValidator) HtmlCompatible() bool {
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
	name    string
	message string
	max     int64
}

// Name of the field
func (c *MaxLengthValidator) Name() string {
	return c.name
}

// SetName of the field
func (c *MaxLengthValidator) SetName(name string) {
	c.name = name
}

// SetMessage set error message
func (c *MaxLengthValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// Validate the value
func (c *MaxLengthValidator) Validate(value any) Error {
	if len([]rune(value.(string))) > int(c.max) {
		return createError(c.name, c.message, fmt.Sprintf("Please shorten %s to %d characters or less", c.name, c.max))
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

// HtmlCompatible for this validator
func (c *MaxLengthValidator) HtmlCompatible() bool {
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
	name     string
	message  string
	min      int64
	optional bool
}

// Name of the field
func (c *MinValidator) Name() string {
	return c.name
}

// SetName of the field
func (c *MinValidator) SetName(name string) {
	c.name = name
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
		return createError(c.name, c.message, fmt.Sprintf("Please increase %s to be %v or more", c.name, c.min))
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

// HtmlCompatible for this validator
func (c *MinValidator) HtmlCompatible() bool {
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
	name    string
	message string
	max     int64
}

// Name of the field
func (c *MaxValidator) Name() string {
	return c.name
}

// SetName of the field
func (c *MaxValidator) SetName(name string) {
	c.name = name
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
		return createError(c.name, c.message, fmt.Sprintf("Please decrease %s to be %v or less", c.name, c.max))
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

// HtmlCompatible for this validator
func (c *MaxValidator) HtmlCompatible() bool {
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
	name     string
	message  string
	re       *regexp.Regexp
	optional bool
}

// Name of the field
func (c *PatternValidator) Name() string {
	return c.name
}

// SetName of the field
func (c *PatternValidator) SetName(name string) {
	c.name = name
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
	str := value.(string)
	if c.optional && str == "" {
		return nil
	}
	if c.re.MatchString(str) {
		return nil
	}
	return createError(c.name, c.message, fmt.Sprintf("Please correct %s into a valid format", c.name))
}

// MarshalJSON for this validator
func (c *PatternValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rule    string `json:"rule"`
		Pattern string `json:"pattern"`
		Message string `json:"message,omitempty"`
	}{"pattern", c.re.String(), c.message})
}

// HtmlCompatible for this validator
func (c *PatternValidator) HtmlCompatible() bool {
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
	name     string
	message  string
	optional bool
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Email field must be a valid email address
func Email() *EmailValidator {
	return &EmailValidator{}
}

// Name of the field
func (c *EmailValidator) Name() string {
	return c.name
}

// SetName of the field
func (c *EmailValidator) SetName(name string) {
	c.name = name
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
	str := value.(string)
	if c.optional && str == "" {
		return nil
	}
	if emailRegex.MatchString(str) {
		return nil
	}
	return createError(c.name, c.message, "Please use a valid email address")
}

// HtmlCompatible for this validator
func (c *EmailValidator) HtmlCompatible() bool {
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
	name    string
	message string
	options []any
}

// Name of the field
func (c *OptionsValidator) Name() string {
	return c.name
}

// SetName of the field
func (c *OptionsValidator) SetName(name string) {
	c.name = name
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
	return createError(c.name, c.message, "Please select one of the valid options")
}

// HtmlCompatible for this validator
func (c *OptionsValidator) HtmlCompatible() bool {
	return false
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
	name    string
	message string
	checker func(string, any) Error
}

// Name of the field
func (c *FieldFuncValidator) Name() string {
	return c.name
}

// SetName of the field
func (c *FieldFuncValidator) SetName(name string) {
	c.name = name
}

// SetMessage set error message
func (c *FieldFuncValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// Validate the value
func (c *FieldFuncValidator) Validate(value any) Error {
	return c.checker(c.name, value)
}

// HtmlCompatible for this validator
func (c *FieldFuncValidator) HtmlCompatible() bool {
	return false
}

// FieldFunc for validating with custom function
func FieldFunc(f func(string, any) Error) Validator {
	return &FieldFuncValidator{
		checker: f,
	}
}

//
// ==================== StructFunc ====================
//

// StructFuncValidator validate struct with custom function. Add to rules with .Struct().
type StructFuncValidator struct {
	name    string
	message string
	checker func(any) Error
}

// Name of the field
func (c *StructFuncValidator) Name() string {
	return c.name
}

// SetName of the field
func (c *StructFuncValidator) SetName(name string) {
	c.name = name
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

// HtmlCompatible for this validator
func (c *StructFuncValidator) HtmlCompatible() bool {
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

func createError(name, custom, fallback string) Error {
	if custom != "" {
		return NewError(custom, name)
	}
	return NewError(fallback, name)
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
