package main

import (
    //"io/ioutil"
    //"os"
    //"sync"

    // TEST DATA
    "image"
    "image/color"
    "image/draw"
    "math/rand"
    // TEST DATA

    "bytes"
    "encoding/json"
    "fmt"
    "image/jpeg"
    "net/http"
    "strconv"
    "strings"
    "text/template"
    //"time"

    "github.com/gorilla/mux"
    //"github.com/pbivrell/Web/photoProject/storage"
    "github.com/disintegration/imaging"
    "github.com/disintegration/imageorient"
    //"github.com/pbivrell/ManagedMap"
)

const (
    DISPLAY_TEMPLATE = "display.tpl"
)

var storage *FileSystemWrapper

func main() {
    _, endpoint, authUrl, client := DoOauth()
    fmt.Println(authUrl)
    r := mux.NewRouter()
    r.HandleFunc(endpoint.Path, endpoint.Handler)
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))).Methods("GET")
    r.HandleFunc("/getImage", GetImage).Methods("GET")
    r.HandleFunc("/testGetImage", TestGetImage).Methods("GET")
    //r.HandleFunc("/configure", Configure).Methods("GET")
    r.HandleFunc("/get", GetConfig).Methods("GET")
    r.HandleFunc("/{path:.*}", DisplayPage).Methods("GET")
    go func(){
        s, _ := NewGoogleDriveStorage(<-client)
        storage, _ = NewFileSystemWrapper(s, "root")
        storage.Crawl()
        fmt.Println("Storage Configured")
    }()
    http.ListenAndServe(":8080", r)
}

/*func SampleJSON(w http.ResponseWriter, r *http.Request) {
    config := Config{
        Title1: "Test",
        Title2: "Page",
        Description: "This is a test page! :)",
        RoutingPage: true,
        //Pictures: []string{"A", "B", "C","D","E","F","G","H","I"},
        Pictures: []string{"Australia", "New Zealand"},
    }
    json.NewEncoder(w).Encode(config)
}*/

func DisplayPage(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path
    if strings.LastIndex(title, "/") == len(title) -1{
        title = title[:len(title)-1]
    }
    title = title[strings.LastIndex(title, "/")+1:]
    if title == "" {
        title = "Home"
    }
    title += " | Paul's Photo Project"
    writeTemplate(w, struct{ Title string}{title}, DISPLAY_TEMPLATE)
}

func GetConfig(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    path := r.FormValue("path")
    //if path == "" {
    //    json.NewEncoder(w).Encode(Config{
    //        ErrorMsg: "Page could not be found!",
    //        ErrorExtra: "No path provided",
    //    })
    //    return
    //}
    data := LoadConfig(path)
    json.NewEncoder(w).Encode(data)
}

type Config struct {
    Title1 string `json:"title1"`
    Title2 string `json:"title2"`
    Description string `json:"description"`
    Pictures []string `json:"pictures"`
    RoutingPage bool `json:"routingPage"`
    ErrorMsg string `json:"error,omitempty"`
    ErrorExtra string `json:"error-extra,omitempty"`
}

func LoadConfig(path string) Config {
    //fmt.Println("Path:",path)
    path = "Original-Photos/" + path
    fmt.Println("Path:",path)
    id := storage.PathToId(path)
    c := Config{}
    if len(id) <= 0 {
        c.ErrorMsg = "Page could not be found!"
        c.ErrorExtra = fmt.Sprintf("No such url [%s]", path)
        return c
    }
    fmt.Println("Id:",id[0])
    files, err := storage.List(id[0])
    if err != nil {
            c.ErrorMsg = "Page could not be found!"
            c.ErrorExtra = fmt.Sprintf("Failed to list directory: %v : %v",path, err)
            return c
    }
    configId := ""
    pictureIds := make([]string, 0)
    c.RoutingPage = false
    //fmt.Println("THINGS:",len(pictureIds))
    for _,v := range files {
        if v.Name == "config.json" {
            configId = v.Id
        }else{
            pictureIds = append(pictureIds,v.Name)
        }
        if isDir, err := storage.IsFolder(v.Id); isDir || err != nil {
            c.RoutingPage = true
        }
    }
    c.Pictures = pictureIds

    if configId != "" {
        data, _ := storage.Get(configId)
        json.NewDecoder(data).Decode(c)
    }
    return c
}

/*func GetPageConfig(url string) PageConfig {
    files, err := storage.List(url + "config.json")
    if err != nil || len(files) != 1 {
        fmt.Println("Error more then one config file for url: " + url)
    }
    data, err := storage.Get(files[0])
    if err != nil {
        fmt.Printf("Error from get config %s: %v\n", url, err)
    }
    var config PageConfig
    json.NewDecoder(data).Decode(&config)
    return config
}

type StorageConfig struct {
    Name  string
    Path  string
    Count int
}

type PageConfig struct {
    Title1      string
    Title2      string
    Description string
    Storages    []StorageConfig
}*/

func writeTemplate(w http.ResponseWriter, data interface{}, templateType string) {
    tmpl, err := template.New(templateType).ParseFiles("./static/templates/" + templateType)
    if err != nil {
        http.Error(w, fmt.Errorf("Failed to process template from file: %v", err).Error(), http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, fmt.Errorf("Failed to use data to create template: %v", err).Error(), http.StatusInternalServerError)
        return
    }
}


func GetImage(w http.ResponseWriter, r *http.Request) {
    root := ""
    t := r.FormValue("root")
    if t == "0" {
        root = "Original-Photos"
    }else if t == "1" {
        root = "Normal-Photos"
    }else if t == "2" {
        root = "Small-Photos"
    }else {
        TestGetImage(w,r)
        return
    }
    url := r.FormValue("url")
    name := r.FormValue("name")
    if url != "/" {
        url = "" + url + "/"
    }
    fmt.Println("P:", root + url + name)
    path := storage.PathToId(root + url + name)
    fmt.Println("Path:",path)
    data, err := storage.Get(path[0])
    img, _, err := imageorient.Decode(data)
    if err != nil {
        fmt.Printf("Err %v\n", err)
        return
    }
    img = imaging.Resize(img, 1500, 0, imaging.Lanczos)
    buffer := new(bytes.Buffer)
    if err := jpeg.Encode(buffer, img, nil); err != nil {
        fmt.Println("unable to encode image.")
    }
    w.Header().Set("Content-Type", "image/jpeg")
    w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
    if _, err := w.Write(buffer.Bytes()); err != nil {
        fmt.Println("unable to write image.")
    }
}

func TestGetImage(w http.ResponseWriter, r *http.Request) {
    buffer := new(bytes.Buffer)
    m := image.NewRGBA(image.Rect(0, 0, 500, 500))
    rc := uint8(rand.Intn(256))
    g := uint8(rand.Intn(256))
    b := uint8(rand.Intn(256))
    blue := color.RGBA{rc, g, b, 255}
    draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

    var img image.Image = m

    if err := jpeg.Encode(buffer, img, nil); err != nil {
        fmt.Println("unable to encode image.")
    }

    w.Header().Set("Content-Type", "image/jpeg")
    w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
    if _, err := w.Write(buffer.Bytes()); err != nil {
        fmt.Println("unable to write image.", err)
    }
}
