package xvalid

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"

	"gopkg.in/guregu/null.v3"
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
func (c *RequiredValidator) Validate(value interface{}) Error {
	v := reflect.ValueOf(value)
	zero := false
	kind := v.Kind()
	if !v.IsValid() {
		zero = true
	} else if v.IsZero() {
		zero = true
	} else if (kind == reflect.Ptr || kind == reflect.Interface) && v.Elem().IsZero() {
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
		Rule string `json:"rule"`
	}{"required"})
}

// HTMLCompatible for this validator
func (c *RequiredValidator) HTMLCompatible() bool {
	return true
}

// Required fields must not be zero
func Required() *RequiredValidator {
	return &RequiredValidator{}
}

//
// ==================== MinStr ====================
//

// MinStrValidator field must have minimum length
type MinStrValidator struct {
	name     string
	message  string
	min      int64
	optional bool
}

// Name of the field
func (c *MinStrValidator) Name() string {
	return c.name
}

// SetName of the field
func (c *MinStrValidator) SetName(name string) {
	c.name = name
}

// SetMessage set error message
func (c *MinStrValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// Optional don't validate if the value is zero
func (c *MinStrValidator) Optional() Validator {
	c.optional = true
	return c
}

// Validate the value
func (c *MinStrValidator) Validate(value interface{}) Error {
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
func (c *MinStrValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rule string `json:"rule"`
		Min  int64  `json:"min"`
	}{"minStr", c.min})
}

// HTMLCompatible for this validator
func (c *MinStrValidator) HTMLCompatible() bool {
	return true
}

// MinStr field must have minimum length
func MinStr(min int64) *MinStrValidator {
	return &MinStrValidator{
		min: min,
	}
}

//
// ==================== MaxStr ====================
//

// MaxStrValidator field have maximum length
type MaxStrValidator struct {
	name    string
	message string
	max     int64
}

// Name of the field
func (c *MaxStrValidator) Name() string {
	return c.name
}

// SetName of the field
func (c *MaxStrValidator) SetName(name string) {
	c.name = name
}

// SetMessage set error message
func (c *MaxStrValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// Validate the value
func (c *MaxStrValidator) Validate(value interface{}) Error {
	if len([]rune(value.(string))) > int(c.max) {
		return createError(c.name, c.message, fmt.Sprintf("Please shorten %s to %d characters or less", c.name, c.max))
	}
	return nil
}

// MarshalJSON for this validator
func (c *MaxStrValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rule string `json:"rule"`
		Max  int64  `json:"max"`
	}{"maxStr", c.max})
}

// HTMLCompatible for this validator
func (c *MaxStrValidator) HTMLCompatible() bool {
	return true
}

// MaxStr field have maximum length
func MaxStr(max int64) *MaxStrValidator {
	return &MaxStrValidator{
		max: max,
	}
}

//
// ==================== MinInt ====================
//

// MinIntValidator field have minimum value
type MinIntValidator struct {
	name     string
	message  string
	min      int64
	optional bool
}

// Name of the field
func (c *MinIntValidator) Name() string {
	return c.name
}

// SetName of the field
func (c *MinIntValidator) SetName(name string) {
	c.name = name
}

// SetMessage set error message
func (c *MinIntValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// Optional don't validate if the value is zero
func (c *MinIntValidator) Optional() Validator {
	c.optional = true
	return c
}

// Validate the value
func (c *MinIntValidator) Validate(value interface{}) Error {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := toInt64(value)
		if c.optional && v == 0 {
			return nil
		}
		if v < c.min {
			return createError(c.name, c.message, fmt.Sprintf("Please increase %s to be %d or more", c.name, c.min))
		}
	default:
		if n, ok := value.(null.Int); ok {
			v := n.Int64
			if c.optional && v == 0 {
				return nil
			}
			if v < c.min {
				return createError(c.name, c.message, fmt.Sprintf("Please increase %s to be %d or more", c.name, c.min))
			}
		} else {
			panic(fmt.Errorf("type not supported: %v", rv.Type()))
		}
	}
	return nil
}

// MarshalJSON for this validator
func (c *MinIntValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rule string `json:"rule"`
		Min  int64  `json:"min"`
	}{"minInt", c.min})
}

// HTMLCompatible for this validator
func (c *MinIntValidator) HTMLCompatible() bool {
	return true
}

// MinInt field have minimum value
func MinInt(min int64) *MinIntValidator {
	return &MinIntValidator{
		min: min,
	}
}

//
// ==================== MaxInt ====================
//

// MaxIntValidator field have maximum value
type MaxIntValidator struct {
	name    string
	message string
	max     int64
}

// Name of the field
func (c *MaxIntValidator) Name() string {
	return c.name
}

// SetName of the field
func (c *MaxIntValidator) SetName(name string) {
	c.name = name
}

// SetMessage set error message
func (c *MaxIntValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// Validate the value
func (c *MaxIntValidator) Validate(value interface{}) Error {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := toInt64(value)
		if v > c.max {
			return createError(c.name, c.message, fmt.Sprintf("Please decrease %s to be %d or less", c.name, c.max))
		}
	default:
		if n, ok := value.(null.Int); ok {
			v := n.Int64
			if v > c.max {
				return createError(c.name, c.message, fmt.Sprintf("Please decrease %s to be %d or less", c.name, c.max))
			}
		} else {
			panic(fmt.Errorf("type not supported: %v", rv.Type()))
		}
	}
	return nil
}

// MarshalJSON for this validator
func (c *MaxIntValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rule string `json:"rule"`
		Max  int64  `json:"max"`
	}{"maxInt", c.max})
}

// HTMLCompatible for this validator
func (c *MaxIntValidator) HTMLCompatible() bool {
	return true
}

// MaxInt field have maximum value
func MaxInt(max int64) *MaxIntValidator {
	return &MaxIntValidator{
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

// Optional don't validate if the value is zero
func (c *PatternValidator) Optional() Validator {
	c.optional = true
	return c
}

// Validate the value
func (c *PatternValidator) Validate(value interface{}) Error {
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
	}{"pattern", c.re.String()})
}

// HTMLCompatible for this validator
func (c *PatternValidator) HTMLCompatible() bool {
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
	PatternValidator
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Email field must be a valid email address
func Email() *EmailValidator {
	return &EmailValidator{
		PatternValidator{
			re:      emailRegex,
			message: "Please use a valid email address",
		},
	}
}

// MarshalJSON for this validator
func (c *EmailValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rule string `json:"rule"`
	}{"email"})
}

// IsEmail returns true if the string is an email
func IsEmail(email string) bool {
	return emailRegex.MatchString(email)
}

//
// ==================== FieldFunc ====================
//

// FieldFuncValidator for validating with custom function
type FieldFuncValidator struct {
	name    string
	message string
	checker func(string, interface{}) Error
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
func (c *FieldFuncValidator) Validate(value interface{}) Error {
	return c.checker(c.name, value)
}

// HTMLCompatible for this validator
func (c *FieldFuncValidator) HTMLCompatible() bool {
	return false
}

// FieldFunc for validating with custom function
func FieldFunc(f func(string, interface{}) Error) Validator {
	return &FieldFuncValidator{
		checker: f,
	}
}

//
// ==================== StructFunc ====================
//

// StructFuncValidator validate struct with custom function
type StructFuncValidator struct {
	name    string
	message string
	checker func(interface{}) Error
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
func (c *StructFuncValidator) Validate(value interface{}) Error {
	return c.checker(value)
}

// HTMLCompatible for this validator
func (c *StructFuncValidator) HTMLCompatible() bool {
	return false
}

// StructFunc validate struct with custom function
func StructFunc(f func(interface{}) Error) Validator {
	return &StructFuncValidator{
		checker: f,
	}
}

func createError(name, custom, fallback string) Error {
	if custom != "" {
		return NewError(custom, name)
	}
	return NewError(fallback, name)
}

func toInt64(value interface{}) int64 {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	}
	panic(fmt.Errorf("cannot convert %v to int64", v.Kind()))
}
