package xvalid

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRequired(t *testing.T) {
	type requiredType struct {
		Number int
		Float  float32
		String string
		Bytes  []byte
		Time   time.Time
		Ptr    *string
	}
	r := requiredType{}
	errs := New(&r).
		Field(&r.Number, Required()).
		Field(&r.Float, Required()).
		Field(&r.String, Required()).
		Field(&r.Bytes, Required()).
		Field(&r.Time, Required()).
		Field(&r.Ptr, Required()).
		Validate(r)
	assert.Len(t, errs.(Errors), 6, "All not set")

	s := ""
	r = requiredType{
		Bytes: make([]byte, 0),
		Time: time.Time{},
		Ptr:  &s,
	}
	errs = New(&r).
		Field(&r.Bytes, Required()).
		Field(&r.Time, Required()).
		Field(&r.Ptr, Required()).
		Validate(r)
	assert.Len(t, errs.(Errors), 3, "Complex type set but empty")

	s = "ok"
	r = requiredType{
		Number: 1,
		Float:  1,
		String: "ok",
		Time:   time.Now(),
		Ptr:    &s,
	}
	errs = New(&r).
		Field(&r.Number, Required()).
		Field(&r.String, Required()).
		Field(&r.Time, Required()).
		Field(&r.Float, Required()).
		Field(&r.Ptr, Required()).
		Validate(r)
	assert.Nil(t, errs, "All value are set")

	msg := "custom message"
	assert.Equal(t, msg, New(&r).Field(&r.Number, Required().SetMessage(msg)).
		Validate(requiredType{}).(Errors)[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&r).Field(&r.Number, Required()).
		Validate(requiredType{}).(Errors)[0].Error(), "Default error message")
}

func TestMinLength(t *testing.T) {
	type strType struct {
		Field string
	}
	str := strType{}
	rules := New(&str).Field(&str.Field, MinLength(2))
	assert.Nil(t, rules.Validate(strType{Field: "123"}), "Long enough")
	assert.Nil(t, rules.Validate(strType{Field: "12"}), "Exactly hit min")
	assert.Len(t, rules.Validate(strType{Field: "1"}).(Errors), 1, "Too short")
	assert.Len(t, rules.Validate(strType{Field: "£"}).(Errors), 1, "Multi-byte characters too short")
	msg := "custom message"
	assert.Equal(t, msg, New(&str).Field(&str.Field, MinLength(2).SetMessage(msg)).
		Validate(strType{Field: "1"}).(Errors)[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&str).Field(&str.Field, MinLength(2)).
		Validate(strType{Field: "1"}).(Errors)[0].Error(), "Default error message")
	// optional
	rules = New(&str).Field(&str.Field, MinLength(3).Optional())
	assert.Nil(t, rules.Validate(strType{Field: ""}), "Invalid but zero")
	assert.Len(t, rules.Validate(strType{Field: " "}).(Errors), 1, "Invalid and not zero")
	assert.Nil(t, rules.Validate(strType{Field: "123"}), "Valid and not zero")
}

func TestMaxLength(t *testing.T) {
	type strType struct {
		Field string
	}
	str := strType{}
	rules := New(&str).Field(&str.Field, MaxLength(2))
	assert.Len(t, rules.Validate(strType{Field: "123"}).(Errors), 1, "Short enough")
	assert.Nil(t, rules.Validate(strType{Field: "12"}), "Exactly hit max")
	assert.Nil(t, rules.Validate(strType{Field: "1"}), "Short enough")
	assert.Nil(t, rules.Validate(strType{Field: "世界"}), "Multi-byte characters are short enough")
	msg := "custom message"
	assert.Equal(t, msg, New(&str).Field(&str.Field, MaxLength(0).SetMessage(msg)).
		Validate(strType{Field: "1"}).(Errors)[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&str).Field(&str.Field, MaxLength(0)).
		Validate(strType{Field: "1"}).(Errors)[0].Error(), "Default error message")
}

func TestMinInt(t *testing.T) {
	type intType struct {
		Int   int
		Float float64
	}
	i := intType{}
	rules := New(&i).Field(&i.Int, Min(0))
	assert.Nil(t, rules.Validate(intType{Int: 1}), "Big enough")
	assert.Nil(t, rules.Validate(intType{Int: 0}), "Exactly hit min")
	assert.Len(t, rules.Validate(intType{Int: -1}).(Errors), 1, "Too low")
	msg := "custom message"
	assert.Equal(t, msg, New(&i).Field(&i.Int, Min(0).SetMessage(msg)).
		Validate(intType{Int: -1}).(Errors)[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&i).Field(&i.Int, Min(0)).
		Validate(intType{Int: -1}).(Errors)[0].Error(), "Default error message")
	// optional
	rules = New(&i).Field(&i.Int, Min(5).Optional())
	assert.Nil(t, rules.Validate(intType{Int: 0}), "Invalid but zero")
	assert.Len(t, rules.Validate(intType{Int: 1}).(Errors), 1, "Invalid and not zero")
	assert.Nil(t, rules.Validate(intType{Int: 5}), "Valid and not zero")

	// float
	rules = New(&i).Field(&i.Float, Min(0))
	assert.Nil(t, rules.Validate(intType{Float: 1}), "Big enough")
	assert.Nil(t, rules.Validate(intType{Float: 0}), "Exactly hit min")
	assert.Len(t, rules.Validate(intType{Float: -1}).(Errors), 1, "Too low")
	assert.Equal(t, msg, New(&i).Field(&i.Float, Min(0).SetMessage(msg)).
		Validate(intType{Float: -1}).(Errors)[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&i).Field(&i.Float, Min(0)).
		Validate(intType{Float: -1}).(Errors)[0].Error(), "Default error message")
	// optional
	rules = New(&i).Field(&i.Float, Min(5).Optional())
	assert.Nil(t, rules.Validate(intType{Float: 0}), "Invalid but zero")
	assert.Len(t, rules.Validate(intType{Float: 1}).(Errors), 1, "Invalid and not zero")
	assert.Nil(t, rules.Validate(intType{Float: 5}), "Valid and not zero")
}

func TestMaxInt(t *testing.T) {
	type intType struct {
		Field int
		Float float64
	}
	i := intType{}
	rules := New(&i).Field(&i.Field, Max(0))
	assert.Len(t, rules.Validate(intType{Field: 1}).(Errors), 1, "Too big")
	assert.Nil(t, rules.Validate(intType{Field: 0}), "Exactly hit max")
	assert.Nil(t, rules.Validate(intType{Field: -1}), "Low engouh")
	msg := "custom message"
	assert.Equal(t, msg, New(&i).Field(&i.Field, Max(0).SetMessage(msg)).
		Validate(intType{Field: 1}).(Errors)[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&i).Field(&i.Field, Max(0)).
		Validate(intType{Field: 1}).(Errors)[0].Error(), "Default error message")

	// float
	rules = New(&i).Field(&i.Float, Max(0))
	assert.Len(t, rules.Validate(intType{Float: 1}).(Errors), 1, "Too big")
	assert.Nil(t, rules.Validate(intType{Float: 0}), "Exactly hit max")
	assert.Nil(t, rules.Validate(intType{Float: -1}), "Low engouh")
	assert.Equal(t, msg, New(&i).Field(&i.Float, Max(0).SetMessage(msg)).
		Validate(intType{Float: 1}).(Errors)[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&i).Field(&i.Float, Max(0)).
		Validate(intType{Float: 1}).(Errors)[0].Error(), "Default error message")
}

func TestPattern(t *testing.T) {
	type patternType struct {
		Field string
	}
	p := patternType{}
	rules := New(&p).Field(&p.Field, Pattern(`\d{2}`))
	assert.Nil(t, rules.Validate(patternType{Field: "00"}), "Exact match")
	assert.Nil(t, rules.Validate(patternType{Field: "1234"}), "Submatch also works")
	assert.Len(t, rules.Validate(patternType{Field: "wrong"}).(Errors), 1, "Pattern is wrong")
	msg := "custom message"
	assert.Equal(t, msg, New(&p).Field(&p.Field, Pattern(`\d{2}`).SetMessage(msg)).
		Validate(patternType{Field: "message"}).(Errors)[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&p).Field(&p.Field, Pattern(`\d{2}`)).
		Validate(patternType{Field: "message"}).(Errors)[0].Error(), "Default error message")
	// optional
	rules = New(&p).Field(&p.Field, Pattern(`\w{3,}`).Optional())
	assert.Nil(t, rules.Validate(patternType{Field: ""}), "Invalid but zero")
	assert.Len(t, rules.Validate(patternType{Field: " "}).(Errors), 1, "Invalid and not zero")
	assert.Nil(t, rules.Validate(patternType{Field: "123"}), "Valid and not zero")
}

func TestEmail(t *testing.T) {
	type emailType struct {
		Field string
	}
	p := emailType{}
	rules := New(&p).Field(&p.Field, Email())
	assert.Len(t, rules.Validate(emailType{Field: "fake"}).(Errors), 1, "Invalid email address")
	assert.Nil(t, rules.Validate(emailType{Field: "test@mail.com"}), "Valid email address")
	msg := "custom message"
	assert.Equal(t, msg, New(&p).Field(&p.Field, Email().SetMessage(msg)).
		Validate(emailType{Field: "invalid"}).(Errors)[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&p).Field(&p.Field, Email()).
		Validate(emailType{Field: "invalid"}).(Errors)[0].Error(), "Default error message")
	// optional
	rules = New(&p).Field(&p.Field, Email().Optional())
	assert.Nil(t, rules.Validate(emailType{Field: ""}), "Invalid but zero")
	assert.Len(t, rules.Validate(emailType{Field: "fake"}).(Errors), 1, "Invalid and not zero")
	assert.Nil(t, rules.Validate(emailType{Field: "test@mail.com"}), "Valid and not zero")
}

func TestFieldFunc(t *testing.T) {
	type funcTest struct {
		Field string
	}
	checker := func(fieldName string, value any) Error {
		if value.(string) == "invalid" {
			return NewError("Invalid field", fieldName)
		}
		return nil
	}
	p := funcTest{}
	rules := New(&p).Field(&p.Field, FieldFunc(checker))
	assert.Nil(t, rules.Validate(funcTest{Field: "valid"}), "Valid")
	assert.Len(t, rules.Validate(funcTest{Field: "invalid"}).(Errors), 1, "Invalid")
}

func TestStructFunc(t *testing.T) {
	type funcTest struct {
		Field string
	}
	checker := func(value any) Error {
		if value.(string) == "invalid" {
			return NewError("Invalid field", "")
		}
		return nil
	}
	p := funcTest{}
	rules := New(&p).Field(&p.Field, StructFunc(checker))
	assert.Nil(t, rules.Validate(funcTest{Field: "valid"}), "Valid")
	assert.Len(t, rules.Validate(funcTest{Field: "invalid"}).(Errors), 1, "Invalid")
}

func TestStruct(t *testing.T) {
	u := structSubject{}
	rules := New(&u).Struct(&compareValidator{})
	assert.Nil(t, rules.Validate(structSubject{Less: 1, More: 2}), "Valid")
	assert.Len(t, rules.Validate(structSubject{Less: 2, More: 1}).(Errors), 1, "Invalid")
}

func TestMarshalJSON(t *testing.T) {
	type strType struct {
		Field string
	}

	str := strType{}
	rules := New(&str).Field(&str.Field, MinStr(2))
	data, _ := json.Marshal(rules.validators)

	assert.Equal(t, string(data), `[{"rule":"minStr","min":2}]`)

	rules = New(&str).Field(&str.Field, MinStr(2).SetMessage("length minimum 2"))
	data, _ = json.Marshal(rules.validators)

	assert.Equal(t, string(data), `[{"rule":"minStr","min":2,"msg":"length minimum 2"}]`)
}

type structSubject struct {
	Less int
	More int
}

type compareValidator struct {
	name    string
	message string
}

func (c *compareValidator) Name() string {
	return c.name
}

func (c *compareValidator) SetName(name string) {
	c.name = name
}

func (c *compareValidator) SetMessage(msg string) Validator {
	c.message = msg
	return c
}

// HtmlCompatible for this validator
func (c *compareValidator) HtmlCompatible() bool {
	return true
}

func (c *compareValidator) Validate(value any) Error {
	subject := value.(structSubject)
	if subject.Less > subject.More {
		return NewError("comparison failed", "")
	}
	return nil
}

func TestJson(t *testing.T) {
	type exportType struct {
		Str   string `json:"string"`
		Int   int    `json:"number"`
		Email string
	}
	e := exportType{}
	rules := New(&e).Field(&e.Str, Required(), MaxLength(5)).Field(&e.Int, Min(10).Optional()).Field(&e.Email, Email())
	j, err := json.MarshalIndent(rules, " ", " ")
	assert.Nil(t, err, "Export rules to json")
	assert.NotEqual(t, len(j), "Export rules to json")
	// println(string(j))
}
