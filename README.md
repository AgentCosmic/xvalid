# xvalid

xvalid is a lightweight validation library that can export rules as JSON so browsers can apply the same rules.

Documentation at [godoc.org](https://godoc.org/github.com/AgentCosmic/xvalid)

## Goals

1. Must be able to export rules so clients can consume them.
2. Only support common rules such as those found in browsers and GUI libraries e.g. length, min, max etc.
3. Must be remain legible as number of rules grows and becomes complex.

## Comparison

Popuplar validation libraries like [go-playground/validate](https://github.com/go-playground/validator) and
[govalidator](https://github.com/asaskevich/govalidator) are great libraries but they suffer from a few problems.
1. Since rules are defined in struct tags, errors are less easy to detect, and it becomes too diffult to read when
there are many rules and many other struct tags defined. By using methods to define rules, we can rely on compilation
checks and can format our code freely.
2. They compile more regex validators than most project will ever need. By
defining only comman validators, we reduce uncessary performance hit; and it is also trivial to copy/paste regex
defined by other libraries.
3. Without being able to export validation rules to client apps, developers will need
to be mindful of keeping the rules in sync. By reusing validation rules, you reduce the chance of client and server
validating wrongly.

## Examples

To define your own validator, you must implement the
[Validator](https://godoc.org/github.com/AgentCosmic/xvalid#Validator) interface. For example, see any of the
validators in [validators.go](https://github.com/AgentCosmic/xvalid/blob/master/validators.go)

Define rules and validate objects:

```go
// Store model
type Store struct {
	Name          string     `json:"name"`
	Address       string     `json:"address"`
	Description   string     `json:"description"`
	Tax           int        `json:"tax"`
}

// ValidationRules for this model.
func (store Store) ValidationRules() xvalid.Rules {
	return xvalid.New(&store).
		Field(&store.Name, xvalid.MinStr(4).Optional().SetMessage("Please lengthen name to 4 characters or more"),
			xvalid.MaxStr(80).SetMessage("Please shorten name to 80 characters or less"),
			xvalid.Pattern("^[a-zA-Z0-9_]+$").Optional().SetMessage("Name may contain alphabets, numbers and underscores"),
			xvalid.Pattern("[a-zA-Z]").Optional().SetMessage("Name must contain at least 1 alphabet"),
			xvalid.FieldFunc(func(fieldName string, value interface{}) xvalid.Error {
				name := value.(string)
				if name == "" {
					return nil
				}
				if name == "admin" {
					return xvalid.NewError("This name is not allowed", fieldName)
				}
				return nil
			})).
		Field(&store.Address, xvalid.Required(), xvalid.MaxStr(120)).
		Field(&store.Description, xvalid.MaxStr(1500)).
		Field(&store.Tax, xvalid.MinInt(0), xvalid.MaxInt(100))
}

// Validate the model and return errors if validation fails.
func (store Store) Validate() error {
	return store.ValidationRules().Validate(store)
}

// then validate
store := Store{}
err := store.Validate()
if err != nil {
	panic(err)
}
```

Export rules as JSON:

```go
func createValidationMap(model ...ruleHolder) map[string]interface{} {
	rules := make(map[string]interface{})
	for _, m := range model {
		rules[reflect.ValueOf(m).Type().Name()] = mapRules(m)
	}
	return rules
}

func mapRules(model ruleHolder) map[string][]interface{} {
	validators := model.ValidationRules().Validators()
	rmap := make(map[string][]interface{})
	for _, v := range validators {
		if !v.HTMLCompatible() {
			continue
		}
		name := v.Name()
		rules, ok := rmap[name]
		if !ok {
			rules = make([]interface{}, 0)
		}
		rules = append(rules, v)
		rmap[name] = rules
	}
	return rmap
}


func main() {
	rules := createValidationMap(model.Member{}, model.Store{}, model.Product{})
	b, err := json.MarshalIndent(rules, "", "	")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
```
