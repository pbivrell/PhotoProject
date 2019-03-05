package main

import (
	"fmt"
	"reflect"
	"unicode"
)

type QueryPart struct {
	Field    string
	Value    interface{}
	Operator string
}

func BuildQuery(qp ...QueryPart) (string, []error) {
    errors := make([]error, len(qp))
    for i, val
}

func EvalQuery(check interface{}, query ...QueryPart) bool {
	reflectedCheck := reflect.ValueOf(check)
	if reflectedCheck.Kind() != reflect.Struct {
		return false
	}
	checkData := make(map[string]interface{})

	for i := 0; i < reflectedCheck.NumField(); i++ {
		if !unicode.IsLower([]rune(reflectedCheck.Type().Field(i).Name)[0]) {
			checkData[reflectedCheck.Type().Field(i).Name] = reflectedCheck.Field(i).Interface()
		}
	}
	
	for _, v := range query {
        evaled := false
		if value, has := checkData[v.Field]; has {
            evaled := value == v.Value
		}

	}
	return true

}

type QueryEvaluator struct {
	Expected string
}

func (q QueryEvaluator) Next() func(string) bool {
	if q.Expected == "operator"
}

func Operator(o string) bool {
	if strings.ToLower(o) == "and" || strings.ToLower(o) == "or" {
		
	}
}

func main() {
	fmt.Println(EvalQuery("apple"))
	fmt.Println(EvalQuery(struct{ A int }{0}, QueryPart{Field: "A", Value: 0, Negate: false, Operator: "and"}))
	fmt.Println(EvalQuery(struct {
		B string
		C bool
	}{"test", false}))
	fmt.Println(EvalQuery(struct{ D struct{ E int } }{D: struct{ E int }{1}}))
}

import (
    "reflect"
)

type QueryPart struct {
    field string
    value string
    negate bool
    operator string
}

type CheckPair struct {
    Type reflect.Type
    Value reflect.Value
}

func EvalQuery(query QueryPart... interface{}) bool {
    reflectedCheck := reflect.ValueOf(check)
    if reflectedCheck.Field(i).Kind() != reflect.Struct {
        return false
    }
    checkData := make(map[string]CheckPair)

    for i := 0; i < reflectedCheck.NumField(); i++ {
        checkData[reflectedCheck.Type().Field(i).Name] = CheckPair{ Type: , Value: 

    fmt.Println(v.Field(i).Interface())
    fmt.Println(v.Type().Field(i).Name)
    fmt.Println(v.Field(i).Type())


    fmt.Println(v.Type().Field(i).Name)
}
}





