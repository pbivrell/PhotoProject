package storage

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





