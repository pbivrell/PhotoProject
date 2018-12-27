package main

import (
    "html/template"
    "net/http"
    "fmt"
    "math"
)

type Todo struct {
    Title string
    Done  bool
}

type TodoPageData struct {
    PageTitle string
    Todos     []Todo
}

type ColContainer struct {
    Max int
    NumMax int
    Next int
}

type Data struct {
    Title1 string
    Title2 string
    Description string
    BigImages []string
    TinyImages []string
}

func main() {
    
    fm := template.FuncMap{"create": func(len int) *ColContainer {
        max := int(math.Ceil(float64(len)/4.0))
        numMax := len % 4
        if numMax == 0{
            numMax = 4
        }
        return &ColContainer{Max: max, NumMax: numMax, Next: max }
    
    }, "newColumn": func(a *ColContainer, i int) bool{
        if i == 0 || i != a.Next {
            return false
        }    
        if (a.Max - 1) == 0 {
            return true
        }
        if i == a.Next{
            if a.NumMax > 1 {
                a.NumMax -= 1
                a.Next += a.Max
            }else{
                a.Next += (a.Max - 1)
            }
            return true
        }
        return false
    }}
    tmpl, err := template.New("template.html").Funcs(fm).ParseFiles("template.html")
    if err != nil{
        fmt.Println("%v\n")
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        d := Data{
            Title1: "Test",
            Title2: "Also a Test",
            Description: "This is a good old description",
            BigImages: []string { "a","b","c","d","e"},
            TinyImages: []string { "a","b","c","d","e"},
        }
        err = tmpl.Execute(w, d)
        if err != nil {
            fmt.Printf("%v",err)
        }
     })

    err = http.ListenAndServe(":8080", nil)
    fmt.Printf("%v\n",err)
}
