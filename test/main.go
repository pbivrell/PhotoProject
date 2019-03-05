package main

import "fmt"

type A interface{
    B() string
}

type Apple struct{
    Value string
}

func (a *Apple) B() string {
    return a.Value
}

func Test(tester A) {
    fmt.Println(tester.B())
}

func main(){
    t := &Apple{ "Hi Paul"}
    Test(t)
    
    var w A
    w = &Apple{ "Second" }
    Test(w)
}
