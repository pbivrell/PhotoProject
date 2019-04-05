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
    //"strings"
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

var storage *GoogleDriveStorage

func main() {
    endpoint, authUrl, client := DoOauth()
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
        storage, _ = NewGoogleDriveStorage(<-client)
        fmt.Println("Storage Configured")
    }()
    http.ListenAndServe(":8080", r)
}

func DisplayPage(w http.ResponseWriter, r *http.Request) {
    writeTemplate(w, nil, DISPLAY_TEMPLATE)
}

func GetConfig(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    path, has := r.URL.Query()["path"]
    if !has {
        json.NewEncoder(w).Encode(Config{
            ErrorMsg: "Page could not be found!",
            ErrorExtra: "No path provided",
        })
        return
    }
    data := LoadConfig(path[0])
    json.NewEncoder(w).Encode(data)
}

type Config struct {
    Title1 string `json:"title1,omitempty"`
    Title2 string `json:"title2,omitempty"`
    Description string `json:"description,omitempty"`
    Pictures []string `json:"pictures,omitempty"`
    RoutingPage bool `json:"routingPage,omitempty"`
    ErrorMsg string `json:"error,omitempty"`
    ErrorExtra string `json:"error-extra,omitempty"`
}

func LoadConfig(path string) Config {
    dir, err := storage.PathSearch(path)
    if err != nil || len(dir) < 1{
        return Config{
            ErrorMsg: "Page could not be found!",
            ErrorExtra: fmt.Sprintf("No such url [%s]: %v", path, err),
        }
    }
    files, err := storage.List(dir[0].Id)
    if err != nil {
        return Config{
            ErrorMsg: "Page could not be found!",
            ErrorExtra: fmt.Sprintf("Failed to list directory: %v", path, err),
        }
    }
    configId := ""
    hasDirectories := false
    pictureIds := make([]string, 0)
    for _,v := range files {
        fmt.Println(v)
        if v.Name == "config.json" {
            configId = v.Id
        }else if isDir, err := storage.IsFolder(v.Id); isDir || err != nil {
            hasDirectories = true
            pictureIds = append(pictureIds,v.Name)
        }else{
            pictureIds = append(pictureIds,v.Id)
        }
    }
    config := Config{
        Title1: "",
        Title2: "",
        Description: "",
        Pictures: pictureIds,
        RoutingPage: hasDirectories,
    }
    if configId != "" {
        data, _ := storage.Get(configId)
        json.NewDecoder(data).Decode(&config)
    }
    return config
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
    var id string
    ids, has := r.URL.Query()["id"]
    if has && len(ids) > 0 {
        id = ids[0]
    }
    data, err := storage.Get(id)
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
    fmt.Println("Called")
    n := rand.Intn(501)
    buffer := new(bytes.Buffer)
    m := image.NewRGBA(image.Rect(0, 0, 500, n+500))
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
