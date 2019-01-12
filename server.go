package photoProject

import (


        //"encoding/json"
        "io/ioutil"
        //"log"
        //"net/http"
        //"os"
        //"image/jpeg"

        "golang.org/x/net/context"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"

        "net/http"
        "fmt"
        "encoding/json"
        //"time"
        "strings"
        "errors"
        "sync" 
       
        "github.com/pbivrell/Web/servable"
        "github.com/gorilla/mux"
        "google.golang.org/api/drive/v3"
       )

const static_dir = "./static/"


// -------------------- Server Construction Functions ----------------

type Server struct{
    GDSCredentialsConfig *oauth2.Config
    GDSClient *drive.Service
    GDSClientLock *sync.RWMutex
}

func NewServer() (*Server) {
    config, err := LoadGDSCredentials()
    if err != nil {
        fmt.Printf("ERROR: Unable to create GDS config form credentials.json: CAUSED BY: %v\n", err)
        return &Server{GDSCredentialsConfig: nil, GDSClient: nil, GDSClientLock: &sync.RWMutex{}}
    }
    return &Server{GDSCredentialsConfig: config, GDSClient: nil, GDSClientLock: &sync.RWMutex{}}
}

const credentialsFile = "./credentials.json"

func LoadGDSCredentials() (*oauth2.Config, error){
        b, err := ioutil.ReadFile(credentialsFile)
        if err != nil {
            return nil, err
        }

        config, err := google.ConfigFromJSON(b, drive.DriveScope)
        if err != nil {
            return nil, err
        }
        return config, nil
}

// -------------------- Implement Servable Interface ---------------

func (s *Server) SubDomains() map[string]func(*mux.Router){
    return map[string]func(*mux.Router){}
}

func (s *Server) ConfigureRouter(r *mux.Router){
    servable.ConfigureSubDomains(s,r)
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(static_dir))))
    r.HandleFunc("/display", s.DisplayPageTemplate)
    r.HandleFunc("/create", s.CreatePageTemplate)
    r.HandleFunc("/load", s.Load)
    r.HandleFunc("/auth", s.AuthGDS)
    r.HandleFunc("/configureGDS", s.ConfigureGDS)
    
    // Error pages
    r.HandleFunc("/notAuthenticated", func (w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w,"Unable to create Google Drive Service. Website admin has not linked authenticated a Google Drive account. Contact admin to fix.")})
    r.HandleFunc("/badCredentials", func (w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w,"Unable to create configuration from Google API Credentials. This webpage will not have access to Google Drive until the server has been restart. See documentation for more information.")})
}

// ---------------------- Server Utiliy Functions -----------

func (s *Server) SetupGDS() (bool, bool) {
    if s.GDSCredentialsConfig == nil {
        return false,false
    }

    s.GDSClientLock.RLock()
    hasClient := s.GDSClient != nil
    s.GDSClientLock.RUnlock()
    return true, hasClient
}

// --------------------- Routes -----------------------

func (s *Server) AuthGDS(w http.ResponseWriter, r *http.Request) {
    if credentials, _ := s.SetupGDS(); !credentials {
        http.Redirect(w,r, "/badCredentials", http.StatusSeeOther)
    }else{
        url := s.GDSCredentialsConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
        http.Redirect(w,r, url, http.StatusSeeOther)
    }
}

func (s *Server) ConfigureGDS(w http.ResponseWriter, r *http.Request) {
    
    // TODO Add security so this endpoint can only be called once and
    // can only be called by google URL
    
    ctx := context.Background()
    code := r.FormValue("code")
    tok, err := s.GDSCredentialsConfig.Exchange(context.Background(), code)
    if err != nil {
        fmt.Fprintf(w, "Failed to create Google Drive Client: Reason: %v\n", err)
        return 
    }
    client := s.GDSCredentialsConfig.Client(ctx, tok)
    srv, err := drive.New(client)
    if err != nil {
        fmt.Fprintf(w, "Failed to create Google Drive Client: Reason: %v\n", err)
        return
    }    
    s.GDSClientLock.Lock()
    s.GDSClient = srv
    s.GDSClientLock.Unlock()
    http.Redirect(w,r, "/create", http.StatusSeeOther)
}

func (s *Server) DisplayPageTemplate(w http.ResponseWriter, r *http.Request) {
    projectID, exists := r.URL.Query()["id"]
    if !exists {
        fmt.Fprintf(w,"Page could not be found")
        return
    }
    resp, err := http.Get("https://drive.google.com/uc?export=view&id=" + projectID[0])
    if err != nil {
        fmt.Fprintf(w, "Page could not be found: %v\n", err)
        return
    }
    decoder := json.NewDecoder(resp.Body)
    var data LoadData
    err = decoder.Decode(&data)
    if err != nil {
        fmt.Fprintf(w, "Page could not be built from source: %v\n", err)
        return
    }
    tmpl, err := getDisplayTemplate()
    if err != nil{
        fmt.Printf("%v\n",err)
        return
    }
    err = tmpl.Execute(w, data)
    if err != nil {
        fmt.Printf("%v\n",err)
    }
}

type LoadData struct{
    Title1 string `json:"title1"`
    Title2 string `json:"title2"`
    Description string `json:"description"`
    Link string `json:"link"`
    folderId string
    BigImages []string `json:"BigImages,omitempty"`
    TinyImages []string `json:"TinyImages,omitempty"`
}

func (s *Server) CreatePageTemplate(w http.ResponseWriter, r *http.Request) {
    tmpl, err := getCreateTemplate()
    if err != nil{
        fmt.Printf("%v\n",err)
    }
    err = tmpl.Execute(w,"")
    if err != nil {
        fmt.Printf("%v\n",err)
    }
}

func (s *Server) Load(w http.ResponseWriter, r *http.Request){
    decoder := json.NewDecoder(r.Body)
    var data LoadData
    err := decoder.Decode(&data)
    if err != nil {
        http.Error(w, fmt.Errorf("Failed to decode JSON: %v\n", err).Error(), 500)
        return
    }
    err = sanatizeInput(&data)
    if err != nil {
        http.Error(w, fmt.Errorf("Failed to sanitize input: %v\n", err).Error(), 500)
        return
    }
    projectId, tiny, big, err := ProcessImages(s.GDSClient, data.folderId)
    if err != nil {
         http.Error(w, err.Error(), 500)
    }
    fmt.Printf("%d %d\n", len(tiny), len(big))
    data.TinyImages = tiny
    data.BigImages = big
    projectId = CreatePageConfig(s.GDSClient, projectId, data, w)
    fmt.Fprintf(w,projectId)
}

func sanatizeInput(data *LoadData) error{
    if uid := strings.Split(data.Link, `folders/`); len(uid) == 2{
        data.folderId= strings.Split(uid[1], `?`)[0]
    }else{
        return errors.New("Invalid Link Address")
    }
    // TODO Limit character lens of title1/2 and description
    // Santaize input... no XSS
    return nil
}

func main(){
    servable.Run(NewServer())
}
