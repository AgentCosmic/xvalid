package xvalid

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

func TestRequired(t *testing.T) {
	type requiredType struct {
		Number int
		String string
		Time   time.Time
		Null   null.Int
		Ptr    *string
	}
	r := requiredType{}
	errs := New(&r).
		Field(&r.Number, Required()).
		Field(&r.String, Required()).
		Field(&r.Time, Required()).
		Field(&r.Null, Required()).
		Field(&r.Ptr, Required()).
		Validate(r)
	assert.Len(t, errs.(Errors), 5, "All not set")

	s := ""
	r = requiredType{
		Null: null.NewInt(0, false),
		Time: time.Time{},
		Ptr:  &s,
	}
	errs = New(&r).
		Field(&r.Time, Required()).
		Field(&r.Null, Required()).
		Field(&r.Ptr, Required()).
		Validate(r)
	assert.Len(t, errs.(Errors), 3, "Complex type set but empty")

	s = "ok"
	r = requiredType{
		Number: 1,
		String: "ok",
		Null:   null.NewInt(0, true),
		Time:   time.Now(),
		Ptr:    &s,
	}
	errs = New(&r).
		Field(&r.Number, Required()).
		Field(&r.String, Required()).
		Field(&r.Time, Required()).
		Field(&r.Null, Required()).
		Field(&r.Ptr, Required()).
		Validate(r)
	assert.Nil(t, errs, "All value are set")

	msg := "custom message"
	assert.Equal(t, msg, New(&r).Field(&r.Number, Required().SetMessage(msg)).
		Validate(requiredType{}).(Errors)[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&r).Field(&r.Number, Required()).
		Validate(requiredType{}).(Errors)[0].Error(), "Default error message")
}

func TestMinStr(t *testing.T) {
	type strType struct {
		Field string
	}
	str := strType{}
	rules := New(&str).Field(&str.Field, MinStr(2))
	assert.Nil(t, rules.Validate(strType{Field: "123"}), "Long enough")
	assert.Nil(t, rules.Validate(strType{Field: "12"}), "Exactly hit min")
	assert.Len(t, rules.Validate(strType{Field: "1"}).(Errors), 1, "Too short")
	assert.Len(t, rules.Validate(strType{Field: "£"}).(Errors), 1, "Multi-byte characters too short")
	msg := "custom message"
	assert.Equal(t, msg, New(&str).Field(&str.Field, MinStr(2).SetMessage(msg)).
		Validate(strType{Field: "1"}).(Errors)[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&str).Field(&str.Field, MinStr(2)).
		Validate(strType{Field: "1"}).(Errors)[0].Error(), "Default error message")
	// optional
	rules = New(&str).Field(&str.Field, MinStr(3).Optional())
	assert.Nil(t, rules.Validate(strType{Field: ""}), "Invalid but zero")
	assert.Len(t, rules.Validate(strType{Field: " "}).(Errors), 1, "Invalid and not zero")
	assert.Nil(t, rules.Validate(strType{Field: "123"}), "Valid and not zero")
}

func TestMaxStr(t *testing.T) {
	type strType struct {
		Field string
	}
	str := strType{}
	rules := New(&str).Field(&str.Field, MaxStr(2))
	assert.Len(t, rules.Validate(strType{Field: "123"}).(Errors), 1, "Short enough")
	assert.Nil(t, rules.Validate(strType{Field: "12"}), "Exactly hit max")
	assert.Nil(t, rules.Validate(strType{Field: "1"}), "Short enough")
	assert.Nil(t, rules.Validate(strType{Field: "世界"}), "Multi-byte characters are short enough")
	msg := "custom message"
	assert.Equal(t, msg, New(&str).Field(&str.Field, MaxStr(0).SetMessage(msg)).
		Validate(strType{Field: "1"}).(Errors)[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&str).Field(&str.Field, MaxStr(0)).
		Validate(strType{Field: "1"}).(Errors)[0].Error(), "Default error message")
}

func TestMinInt(t *testing.T) {
	type intType struct {
		Field int
		Null  null.Int
	}
	i := intType{}
	rules := New(&i).Field(&i.Field, MinInt(0))
	assert.Nil(t, rules.Validate(intType{Field: 1}), "Big enough")
	assert.Nil(t, rules.Validate(intType{Field: 0}), "Exactly hit min")
	assert.Len(t, rules.Validate(intType{Field: -1}).(Errors), 1, "Too low")
	msg := "custom message"
	assert.Equal(t, msg, New(&i).Field(&i.Field, MinInt(0).SetMessage(msg)).
		Validate(intType{Field: -1}).(Errors)[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&i).Field(&i.Field, MinInt(0)).
		Validate(intType{Field: -1}).(Errors)[0].Error(), "Default error message")
	// optional
	rules = New(&i).Field(&i.Field, MinInt(5).Optional())
	assert.Nil(t, rules.Validate(intType{Field: 0}), "Invalid but zero")
	assert.Len(t, rules.Validate(intType{Field: 1}).(Errors), 1, "Invalid and not zero")
	assert.Nil(t, rules.Validate(intType{Field: 5}), "Valid and not zero")

	// null
	rules = New(&i).Field(&i.Null, MinInt(0))
	assert.Nil(t, rules.Validate(intType{Null: null.IntFrom(1)}), "Big enough")
	assert.Nil(t, rules.Validate(intType{Null: null.IntFrom(0)}), "Exactly hit min")
	assert.Len(t, rules.Validate(intType{Null: null.IntFrom(-1)}).(Errors), 1, "Too low")
	assert.Equal(t, msg, New(&i).Field(&i.Null, MinInt(0).SetMessage(msg)).
		Validate(intType{Null: null.IntFrom(-1)}).(Errors)[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&i).Field(&i.Null, MinInt(0)).
		Validate(intType{Null: null.IntFrom(-1)}).(Errors)[0].Error(), "Default error message")
	// optional
	rules = New(&i).Field(&i.Null, MinInt(5).Optional())
	assert.Nil(t, rules.Validate(intType{Null: null.IntFrom(0)}), "Invalid but zero")
	assert.Len(t, rules.Validate(intType{Null: null.IntFrom(1)}).(Errors), 1, "Invalid and not zero")
	assert.Nil(t, rules.Validate(intType{Null: null.IntFrom(5)}), "Valid and not zero")
}

func TestMaxInt(t *testing.T) {
	type intType struct {
		Field int
		Null  null.Int
	}
	i := intType{}
	rules := New(&i).Field(&i.Field, MaxInt(0))
	assert.Len(t, rules.Validate(intType{Field: 1}).(Errors), 1, "Too big")
	assert.Nil(t, rules.Validate(intType{Field: 0}), "Exactly hit max")
	assert.Nil(t, rules.Validate(intType{Field: -1}), "Low engouh")
	msg := "custom message"
	assert.Equal(t, msg, New(&i).Field(&i.Field, MaxInt(0).SetMessage(msg)).
		Validate(intType{Field: 1}).(Errors)[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&i).Field(&i.Field, MaxInt(0)).
		Validate(intType{Field: 1}).(Errors)[0].Error(), "Default error message")

	// null
	rules = New(&i).Field(&i.Null, MaxInt(0))
	assert.Len(t, rules.Validate(intType{Null: null.IntFrom(1)}).(Errors), 1, "Too big")
	assert.Nil(t, rules.Validate(intType{Null: null.IntFrom(0)}), "Exactly hit max")
	assert.Nil(t, rules.Validate(intType{Null: null.IntFrom(-1)}), "Low engouh")
	assert.Equal(t, msg, New(&i).Field(&i.Null, MaxInt(0).SetMessage(msg)).
		Validate(intType{Null: null.IntFrom(1)}).(Errors)[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&i).Field(&i.Null, MaxInt(0)).
		Validate(intType{Null: null.IntFrom(1)}).(Errors)[0].Error(), "Default error message")
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
	checker := func(fieldName string, value interface{}) Error {
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
	checker := func(value interface{}) Error {
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

// HTMLCompatible for this validator
func (c *compareValidator) HTMLCompatible() bool {
	return true
}

func (c *compareValidator) Validate(value interface{}) Error {
	subject := value.(structSubject)
	if subject.Less > subject.More {
		return NewError("comparison failed", "")
	}
	return nil
}
