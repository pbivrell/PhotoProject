package main

import "os"
import "text/template"
import "math"

//{{range $i, $v := .}}{{{$i}} {{$v}}{{end}}
//`

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

//    t := `{{$colContainer:= create (len .)}}{{range $i, $v := .}}{{if newColumn $colContainer $i}} {{end}}{{$v}}{{end}}`
    
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

    d := Data{
        Title1: "Test",
        Title2: "Also a Test",
        Description: "This is a good old description",
        BigImages: []string { "a","b","c","d","e"},
        TinyImages: []string { "a","b","c","d","e"},
    }
    //d := []string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q"}
    template.Must(template.New("").Funcs(fm).ParseFiles("./template.tpl")).Execute(os.Stdout, d)
}
