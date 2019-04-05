package main

import (


	// TEST DATA
	"image"
	"image/color"
	"image/draw"
    "math/rand"
	// TEST DATA

    "strconv"
    "net/http"
    "time"
    "text/template"
    "fmt"
    "strings"
    "bytes"
	"image/jpeg"
    "encoding/json"

    "github.com/gorilla/mux"
    //"github.com/pbivrell/Web/photoProject/storage"
    "github.com/pbivrell/ManagedMap"
    "github.com/disintegration/imaging"
)

const (
    // Templates
    UPLOAD_TEMPLATE = "upload.tpl"
    LOADING_TEMPLATE = "loading.tpl"
    DISPLAY_TEMPLATE = "display.tpl"
    CONFIGURE_TEMPLATE = "configure.tpl"
    ADMIT_TEMPLATE = "admit.tpl"
    NO_ADMIT_TEMPLATE = "noAdmit.tpl"
)

func configureRouter(r *mux.Router) *mux.Router {
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(server.StaticDir)))).Methods("GET")
    // Server HTML webpages
    r.HandleFunc("/upload", UploadPage).Methods("GET")
    r.HandleFunc("/loading", LoadingPage).Methods("GET")
    r.HandleFunc("/configure", ConfigurePage).Methods("GET")
    r.HandleFunc("/admit", AdmitPage).Methods("GET")
    // Preform Functionality
    r.HandleFunc("/uploadProject", UploadProject).Methods("POST")
    r.HandleFunc("/uploadImages", UploadImages).Methods("POST")
    r.HandleFunc("/getImage", GetImage).Methods("GET")
    r.HandleFunc("/testGetImage", TestGetImage).Methods("GET")
    // Status
    r.HandleFunc("/progress", Progress).Methods("GET")
    // This needs to be last because it matches all of the above paths
    r.HandleFunc("/{path:.*}", DisplayPage).Methods("GET")
    // Middleware
    r.Use(middleware)
    return r
}

// Server HTML webpages
func UploadPage(w http.ResponseWriter, r *http.Request) {
    if server.AdmitHash == "" {
        writeTemplate(w, nil, UPLOAD_TEMPLATE)
        return
    }
    token, has := r.URL.Query()["token"]
    if !has {
        http.Error(w, "Unable to access upload page: No token provided\n", http.StatusUnauthorized)
        return
    }
    if _, has := server.TokenMap.Get(token[0]); !has {
        http.Error(w, "Unable to access upload page: Invalid or expired token\n", http.StatusUnauthorized)
        return
    }
    writeTemplate(w, nil, UPLOAD_TEMPLATE)
}

func LoadingPage(w http.ResponseWriter, r *http.Request) {
    writeTemplate(w, nil, LOADING_TEMPLATE)
}

/*func PasswordProtect(r *http.Request) bool {
    if server.AdmitHash == "" {
        return 
    }
    code, has := r.URL.Query()["code"]
    if !has {
        http.Error(w, "Unable to create access token for write storage: Password not provided.", http.StatusUnauthorized)
        return
    }

}*/

func ConfigurePage(w http.ResponseWriter, r *http.Request) {
    setup := make([]string, 0)
    for _, v := range server.WriteStorage {
        setup = append(setup, v.HowToConfigure())
    }
    content := struct{ Setup []string }{ setup }
    writeTemplate(w, content, CONFIGURE_TEMPLATE)
}

func DisplayPage(w http.ResponseWriter, r *http.Request) {
    // TODO WRITE GETTING CONFIGURATION PAGES
    // This should add to a list of storage volumes.  
    //
    config := GetPageConfig(r.URL.Path)
    storageJSON, err := json.Marshal(config.Storages)
    if err != nil {
        fmt.Println("%v\n", err)
    }
    content  := struct{ Title1 string; Title2 string; Description string; StorageJSON string }{ config.Title1, config.Title2, config.Description, string(storageJSON)}
    writeTemplate(w, content, DISPLAY_TEMPLATE)
}

type StorageConfig struct {
    Name string
    Path string
    Count int
}

type PageConfig struct {
    Title1 string
    Title2 string
    Description string
    Storages []StorageConfig
}

func ExposePageConfig(w http.ResponseWriter, r *http.Request) {
    // TODO WORK ON THIS MAYBE
    //w.Header().Set("Content-Type", "application/json")
    //json.NewEncoder(w).Encode(GetPageConfig(w.)
}

func GetPageConfig(url string) PageConfig {
    // TODO GetConfigs from Write Storage
    return PageConfig{ Title1: "Hello", Title2: "There", Description: "This is a test page...", Storages: []StorageConfig{} }
}

func AdmitPage(w http.ResponseWriter, r *http.Request) {
    if server.AdmitHash == "" {
        writeTemplate(w, nil, NO_ADMIT_TEMPLATE)
        return
    }
    code, has := r.URL.Query()["code"]
    if !has {
        http.Error(w, "Unable to create access token for write storage: Password not provided.", http.StatusUnauthorized)
        return
    }
    if CheckPasswordHash(code[0], server.AdmitHash){
        conf := ManagedMap.Config{ManagedMap.DefaultTimeout, ManagedMap.DefaultAccessCount}
        expiration, has := r.URL.Query()["exp"]
        if has && len(expiration) > 0 {
            timeout, err := strconv.ParseInt(expiration[0], 10, 64)
            if err == nil {
                conf.Timeout = time.Duration(timeout) * time.Second
            }
        }
        accessCount, has := r.URL.Query()["uses"]
        if has && len(accessCount) > 0 {
            count, err := strconv.ParseUint(accessCount[0], 10, 64)
            if err == nil {
                conf.AccessCount = count
            }
        }
        token := GetRandomToken()
        for server.TokenMap.Has(token) {
            token = GetRandomToken()
        }
        server.TokenMap.Put(token, nil)
        content  := struct{ Link string; Access uint64; Timeout time.Duration}{ server.IP + ":" + server.Port + "/upload?token=" + token, conf.AccessCount, conf.Timeout }
        writeTemplate(w, content, ADMIT_TEMPLATE)
        return
    }
    http.Error(w, "Unable to create access token for write storage: Incorrect password.", http.StatusUnauthorized)
}

// Preform Functionality
func UploadProject(w http.ResponseWriter, r *http.Request) {

}

func UploadImages(w http.ResponseWriter, r *http.Request) {

}

func GetImage(w http.ResponseWriter, r *http.Request){
    var path string
    paths, has := r.URL.Query()["path"]
    if has && len(paths) > 0 {
        path = paths[0]
    }
    files, err := server.ReadStorage.Search(Query{Name: path})
    if err != nil || len(files) < 1 {
        fmt.Printf("Err %v\n", err)
        return
    }
    data, err := server.ReadStorage.Get(files[0])
    if err != nil {
        fmt.Printf("Err %v\n", err)
        return
    }
    img,_, err := image.Decode(data)
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
    n := rand.Intn(501)
    buffer := new(bytes.Buffer)
    m := image.NewRGBA(image.Rect(0, 0, 500, n + 500))
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
		fmt.Println("unable to write image.")
	}
    
}

// Status
func Progress(w http.ResponseWriter, r *http.Request) {
    fid, err := server.WriteStorage[0].NewFolder("apple","")
    if err != nil {
        fmt.Printf("Error:%v\n", err)
        return
    }else{
        fmt.Printf("FId:%s\n", fid)
    }

    id, err := server.WriteStorage[0].NewFile("test.json", fid, MIME_TYPE_JSON, strings.NewReader("{ \"test\": \"apple\" }"))
    if err != nil {
        fmt.Printf("Error:%v\n", err)
        return
    }else{
        fmt.Printf("Id:%s\n", id)
    }
    err = server.WriteStorage[0].Update(id, strings.NewReader("{ \"changed\"}"))
    if err != nil {
        fmt.Printf("Error:%v\n", err)
        return
    }

    contentReader, err := server.WriteStorage[0].Get(fid)
    if err != nil {
        fmt.Printf("Error:%v\n", err)
        return
    }
    buf := new(bytes.Buffer)
    buf.ReadFrom(contentReader)
    s := buf.String()
    contentReader.Close()
    fmt.Printf("Data: %s\n", s)

    files, err := server.WriteStorage[0].List(fid)
    if err != nil {
        fmt.Printf("Error:%v\n", err)
        return
    }
    for _, f := range files {
        fmt.Printf("Id: %s\n", f)
    }

}

// Middlewear
// TODO Logging Middleware?
func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if strings.Contains(r.URL.Path, "/configure") || server.Configured() {
            next.ServeHTTP(w, r)
            return
        }
        //next.ServeHTTP(w, r)
        http.Redirect(w, r, "/configure", http.StatusSeeOther)
    })
}


func writeTemplate(w http.ResponseWriter, data interface{}, templateType string) {
    tmpl, err := template.New(templateType).ParseFiles(server.StaticDir+ "/templates/" + templateType)
    if err != nil {
        http.Error(w, fmt.Errorf("Failed to process template from file: %v",err).Error(), http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, fmt.Errorf("Failed to use data to create template: %v",err).Error(), http.StatusInternalServerError)
        return
    }
}
