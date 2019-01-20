package photoProject

import (
    "html/template"
    "math"
)

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

const template_dir = "./photoProject/static/templates/"

func template_func_create(len int) *ColContainer {
    max := int(math.Ceil(float64(len)/4.0))
    numMax := len % 4
    if numMax == 0{
        numMax = 4
    }
    return &ColContainer{Max: max, NumMax: numMax, Next: max }
}

func template_func_newColumn(a *ColContainer, i int) bool{
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
}

const displayTemplate = "photoPage.tpl"
func getDisplayTemplate() (*template.Template,error) {
    return  template.New(displayTemplate).Funcs(template.FuncMap{"create":template_func_create, "newColumn":template_func_newColumn,}).ParseFiles(template_dir + displayTemplate)
}

const indexTemplate = "indexPage.tpl"
func getIndexTemplate() (*template.Template, error) {
    return template.New(indexTemplate).ParseFiles(template_dir + indexTemplate)
}

const createTemplate = "createPage.html"
func getCreateTemplate() (*template.Template, error) {
    return template.New(createTemplate).ParseFiles(template_dir + createTemplate)
}

const authTemplate = "authPage.tpl"
func getAuthTemplate() (*template.Template, error) {
    return template.New(authTemplate).ParseFiles(template_dir + authTemplate)
}


