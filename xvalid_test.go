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
	assert.Len(t, errs, 6, "All not set")

	s := ""
	r = requiredType{
		Bytes: make([]byte, 0),
		Time:  time.Time{},
		Ptr:   &s,
	}
	errs = New(&r).
		Field(&r.Bytes, Required()).
		Field(&r.Time, Required()).
		Field(&r.Ptr, Required()).
		Validate(r)
	assert.Len(t, errs, 3, "Complex type set but empty")

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
		Validate(requiredType{})[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&r).Field(&r.Number, Required()).
		Validate(requiredType{})[0].Error(), "Default error message")
}

func TestMinLength(t *testing.T) {
	type strType struct {
		Field string
	}
	str := strType{}
	rules := New(&str).Field(&str.Field, MinLength(2))
	assert.Nil(t, rules.Validate(strType{Field: "123"}), "Long enough")
	assert.Nil(t, rules.Validate(strType{Field: "12"}), "Exactly hit min")
	assert.Len(t, rules.Validate(strType{Field: "1"}), 1, "Too short")
	assert.Len(t, rules.Validate(strType{Field: "£"}), 1, "Multi-byte characters too short")
	msg := "custom message"
	assert.Equal(t, msg, New(&str).Field(&str.Field, MinLength(2).SetMessage(msg)).
		Validate(strType{Field: "1"})[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&str).Field(&str.Field, MinLength(2)).
		Validate(strType{Field: "1"})[0].Error(), "Default error message")
	// optional
	rules = New(&str).Field(&str.Field, MinLength(3).SetOptional())
	assert.Nil(t, rules.Validate(strType{Field: ""}), "Invalid but zero")
	assert.Len(t, rules.Validate(strType{Field: " "}), 1, "Invalid and not zero")
	assert.Nil(t, rules.Validate(strType{Field: "123"}), "Valid and not zero")
}

func TestMaxLength(t *testing.T) {
	type strType struct {
		Field string
	}
	str := strType{}
	rules := New(&str).Field(&str.Field, MaxLength(2))
	assert.Len(t, rules.Validate(strType{Field: "123"}), 1, "Short enough")
	assert.Nil(t, rules.Validate(strType{Field: "12"}), "Exactly hit max")
	assert.Nil(t, rules.Validate(strType{Field: "1"}), "Short enough")
	assert.Nil(t, rules.Validate(strType{Field: "世界"}), "Multi-byte characters are short enough")
	msg := "custom message"
	assert.Equal(t, msg, New(&str).Field(&str.Field, MaxLength(0).SetMessage(msg)).
		Validate(strType{Field: "1"})[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&str).Field(&str.Field, MaxLength(0)).
		Validate(strType{Field: "1"})[0].Error(), "Default error message")
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
	assert.Len(t, rules.Validate(intType{Int: -1}), 1, "Too low")
	msg := "custom message"
	assert.Equal(t, msg, New(&i).Field(&i.Int, Min(0).SetMessage(msg)).
		Validate(intType{Int: -1})[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&i).Field(&i.Int, Min(0)).
		Validate(intType{Int: -1})[0].Error(), "Default error message")
	// optional
	rules = New(&i).Field(&i.Int, Min(5).SetOptional())
	assert.Nil(t, rules.Validate(intType{Int: 0}), "Invalid but zero")
	assert.Len(t, rules.Validate(intType{Int: 1}), 1, "Invalid and not zero")
	assert.Nil(t, rules.Validate(intType{Int: 5}), "Valid and not zero")

	// float
	rules = New(&i).Field(&i.Float, Min(0))
	assert.Nil(t, rules.Validate(intType{Float: 1}), "Big enough")
	assert.Nil(t, rules.Validate(intType{Float: 0}), "Exactly hit min")
	assert.Len(t, rules.Validate(intType{Float: -1}), 1, "Too low")
	assert.Equal(t, msg, New(&i).Field(&i.Float, Min(0).SetMessage(msg)).
		Validate(intType{Float: -1})[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&i).Field(&i.Float, Min(0)).
		Validate(intType{Float: -1})[0].Error(), "Default error message")
	// optional
	rules = New(&i).Field(&i.Float, Min(5).SetOptional())
	assert.Nil(t, rules.Validate(intType{Float: 0}), "Invalid but zero")
	assert.Len(t, rules.Validate(intType{Float: 1}), 1, "Invalid and not zero")
	assert.Nil(t, rules.Validate(intType{Float: 5}), "Valid and not zero")
}

func TestMaxInt(t *testing.T) {
	type intType struct {
		Field int
		Float float64
	}
	i := intType{}
	rules := New(&i).Field(&i.Field, Max(0))
	assert.Len(t, rules.Validate(intType{Field: 1}), 1, "Too big")
	assert.Nil(t, rules.Validate(intType{Field: 0}), "Exactly hit max")
	assert.Nil(t, rules.Validate(intType{Field: -1}), "Low engouh")
	msg := "custom message"
	assert.Equal(t, msg, New(&i).Field(&i.Field, Max(0).SetMessage(msg)).
		Validate(intType{Field: 1})[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&i).Field(&i.Field, Max(0)).
		Validate(intType{Field: 1})[0].Error(), "Default error message")

	// float
	rules = New(&i).Field(&i.Float, Max(0))
	assert.Len(t, rules.Validate(intType{Float: 1}), 1, "Too big")
	assert.Nil(t, rules.Validate(intType{Float: 0}), "Exactly hit max")
	assert.Nil(t, rules.Validate(intType{Float: -1}), "Low engouh")
	assert.Equal(t, msg, New(&i).Field(&i.Float, Max(0).SetMessage(msg)).
		Validate(intType{Float: 1})[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&i).Field(&i.Float, Max(0)).
		Validate(intType{Float: 1})[0].Error(), "Default error message")
}

func TestPattern(t *testing.T) {
	type patternType struct {
		Field string
	}
	p := patternType{}
	rules := New(&p).Field(&p.Field, Pattern(`\d{2}`))
	assert.Nil(t, rules.Validate(patternType{Field: "00"}), "Exact match")
	assert.Nil(t, rules.Validate(patternType{Field: "1234"}), "Submatch also works")
	assert.Len(t, rules.Validate(patternType{Field: "wrong"}), 1, "Pattern is wrong")
	msg := "custom message"
	assert.Equal(t, msg, New(&p).Field(&p.Field, Pattern(`\d{2}`).SetMessage(msg)).
		Validate(patternType{Field: "message"})[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&p).Field(&p.Field, Pattern(`\d{2}`)).
		Validate(patternType{Field: "message"})[0].Error(), "Default error message")
	// optional
	rules = New(&p).Field(&p.Field, Pattern(`\w{3,}`).SetOptional())
	assert.Nil(t, rules.Validate(patternType{Field: ""}), "Invalid but zero")
	assert.Len(t, rules.Validate(patternType{Field: " "}), 1, "Invalid and not zero")
	assert.Nil(t, rules.Validate(patternType{Field: "123"}), "Valid and not zero")
}

func TestEmail(t *testing.T) {
	type emailType struct {
		Field string
	}
	p := emailType{}
	rules := New(&p).Field(&p.Field, Email())
	assert.Len(t, rules.Validate(emailType{Field: "fake"}), 1, "Invalid email address")
	assert.Nil(t, rules.Validate(emailType{Field: "test@mail.com"}), "Valid email address")
	msg := "custom message"
	assert.Equal(t, msg, New(&p).Field(&p.Field, Email().SetMessage(msg)).
		Validate(emailType{Field: "invalid"})[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&p).Field(&p.Field, Email()).
		Validate(emailType{Field: "invalid"})[0].Error(), "Default error message")
	// optional
	rules = New(&p).Field(&p.Field, Email().SetOptional())
	assert.Nil(t, rules.Validate(emailType{Field: ""}), "Invalid but zero")
	assert.Len(t, rules.Validate(emailType{Field: "fake"}), 1, "Invalid and not zero")
	assert.Nil(t, rules.Validate(emailType{Field: "test@mail.com"}), "Valid and not zero")
}

func TestOptions(t *testing.T) {
	type optionsType struct {
		Str string
		Int int
	}
	o := optionsType{}
	// string
	rules := New(&o).Field(&o.Str, Options("a", "b", "c"))
	assert.Len(t, rules.Validate(optionsType{Str: "x"}), 1, "Not in options")
	assert.Nil(t, rules.Validate(optionsType{Str: "b"}), "Valid option")
	msg := "custom message"
	assert.Equal(t, msg, New(&o).Field(&o.Str, Options().SetMessage(msg)).
		Validate(optionsType{Str: "invalid"})[0].Error(), "Custom error message")
	assert.NotEqual(t, msg, New(&o).Field(&o.Str, Options()).
		Validate(optionsType{Str: "invalid"})[0].Error(), "Default error message")
	// int
	rules = New(&o).Field(&o.Int, Options(1, 2, 3))
	assert.Len(t, rules.Validate(optionsType{Int: 5}), 1, "Not in options")
	assert.Nil(t, rules.Validate(optionsType{Int: 1}), "Valid option")
	// mixed type
	rules = New(&o).Field(&o.Int, Options("a", 5, make([]byte, 0)))
	assert.Len(t, rules.Validate(optionsType{Int: -1}), 1, "Not in options")
	assert.Nil(t, rules.Validate(optionsType{Int: 5}), "Valid option")
}

func TestFieldFunc(t *testing.T) {
	type funcTest struct {
		Field string
	}
	checker := func(field []string, value any) Error {
		if value.(string) == "invalid" {
			return NewError("Invalid field", field...)
		}
		return nil
	}
	p := funcTest{}
	rules := New(&p).Field(&p.Field, FieldFunc(checker))
	assert.Nil(t, rules.Validate(funcTest{Field: "valid"}), "Valid")
	assert.Len(t, rules.Validate(funcTest{Field: "invalid"}), 1, "Invalid")
}

func TestStructFunc(t *testing.T) {
	type funcTest struct {
		A int
		B int
	}
	checker := func(value any) Error {
		s := value.(funcTest)
		if s.A > s.B {
			return NewError("custom error")
		}
		return nil
	}
	p := funcTest{}
	rules := New(&p).Struct(StructFunc(checker))
	assert.Nil(t, rules.Validate(funcTest{A: 3, B: 10}), "Valid")
	errs := rules.Validate(funcTest{A: 3, B: 1})
	assert.Len(t, errs, 1, "Invalid")
	assert.Equal(t, errs[0].Error(), "custom error", "Error message")
}

func TestEmbeded(t *testing.T) {
	type Deep struct {
		DeepInt int `json:"deepInt"`
	}
	type Embed struct {
		EmbedStr   string  `json:"embedStr"`
		EmbedFloat float32 `json:"embedFloat"`
		Deep       `json:"deep"`
	}
	type nestedType struct {
		Top string
		Embed
	}
	n := nestedType{}

	rules := New(&n).Field(&n.Top, MinLength(3)).
		Field(&n.Embed.EmbedStr, MaxLength(2), MinLength(1)).
		Field(&n.EmbedFloat, Min(2)).
		Field(&n.DeepInt, Min(5))
	assert.Len(t, rules.Validate(nestedType{Top: "x"}), 4, "All fail")
	assert.Nil(t, rules.Validate(nestedType{Top: "abc", Embed: Embed{EmbedStr: "x", EmbedFloat: 3, Deep: Deep{5}}}), "All pass")
}

func TestMarshalJSON(t *testing.T) {
	type Embed struct {
		EmbedStr string `json:"embedStr"`
	}
	type exportType struct {
		Embed
		Str string
		Int int `json:"number,omitempty"`
	}
	e := exportType{}
	// json rules
	rules := New(&e).
		Field(&e.Str, Required(), MaxLength(5)).
		Field(&e.Int, Min(10).SetOptional().SetMessage("my message")).
		Field(&e.EmbedStr, Required())
	j, _ := json.Marshal(rules)
	assert.Equal(t,
		`{"Str":[{"rule":"required"},{"rule":"maxLength","max":5}],"embedStr":[{"rule":"required"}],"number":[{"rule":"min","min":10,"message":"my message"}]}`,
		string(j), "Export rules to json")
	// json errors
	errs := rules.Validate(e)
	j, _ = json.Marshal(errs)
	assert.Equal(t,
		`[{"message":"Please enter the Str","field":"Str"},{"message":"Please enter the embedStr","field":"embedStr"}]`,
		string(j), "Export errors json")
	// as map
	errsMap := errs.ToMap()
	j, _ = json.Marshal(errsMap)
	assert.Equal(t,
		`{"Str":"Please enter the Str","embedStr":"Please enter the embedStr"}`,
		string(j), "Export errors json as map")
}
