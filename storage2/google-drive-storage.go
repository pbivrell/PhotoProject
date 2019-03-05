package main


import (
    "google.golang.org/api/drive/v3"
    "golang.org/x/net/context"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"

    "github.com/gorilla/mux"

    "fmt"
    "net/http"
    "sync"
    "io/ioutil"
    "io"
)

const (
    STORAGE_TYPE_GOOGLE_DRIVE = "googleDrive"
    MIME_TYPE_GOOGLE_DRIVE_FOLDER = "application/vnd.google-apps.folder"
    MIME_TYPE_JSON = "application/json"
    MIME_TYPE_JPEG= "image/jpeg"
)

type Storage interface {
    // Configuration Methods
    ConfigureEndpoints(r *mux.Router)
    HowToConfigure() string
    Configured() bool
    Configure(interface{}) error
    // Storage Interaction Methods
    NewFolder(name string, parentId string) (string, error)
    NewFile(name string, parentId string, ctype string, content io.Reader) (string, error)
    Update(id string, content io.Reader) error
    Delete(id string) error
    Get(id string) (io.ReadCloser, error)
    List(parentId string) ([]string, error)
    Search(query Query) ([]string, error)
}

type GoogleDriveStorage struct {
    VolumeId string
    Config *oauth2.Config
    ConfigPath string
    Service *drive.Service
    Lock *sync.RWMutex
    RootPath string
}

func NewStorage(storageType string, name string, extra string) *GoogleDriveStorage {
    return NewGoogleDriveStorage(name, extra)
}

func NewGoogleDriveStorage(volumeId string, credentialsFilePath string) *GoogleDriveStorage {
    return &GoogleDriveStorage{
            VolumeId: volumeId,
            Config: nil,
            ConfigPath: credentialsFilePath,
            Service: nil,
            Lock: &sync.RWMutex{},
    }
}

func (s *GoogleDriveStorage) HowToConfigure() string {
    message := fmt.Sprintf("Storage: Google Drive Service \"%s\": %%s", s.VolumeId)
    if s.Configured(){
       return fmt.Sprintf(message, "Configured!")
    }
    s.Lock.Lock()
    defer s.Lock.Unlock()
    if s.Config == nil {
        err := s.loadCredentials()
        if err != nil {
            reason := fmt.Sprintf("Error: Could not load oauth2 credentials file: %s", err)
            return fmt.Sprintf(message, reason)
        }
    }
    url := s.Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
    reason := fmt.Sprintf("Not Configured: Use this <a href=\"%s\">link to configure</a>", url)
    return fmt.Sprintf(message, reason)
}

func (s *GoogleDriveStorage) Configured() bool {
    s.Lock.RLock()
    defer s.Lock.RUnlock()
    return s.Service != nil
}

func (s *GoogleDriveStorage) ConfigureEndpoints(r *mux.Router) {
    s.Lock.RLock()
    defer s.Lock.RUnlock()
    r.HandleFunc("/configureGDS-" + s.VolumeId, s.configureEndpoint)
}

func (s *GoogleDriveStorage) configureEndpoint(w http.ResponseWriter, r *http.Request){
    fmt.Println("Configuring")
    // Retrieve Code from google oauth callback
    code := r.FormValue("code")
    if err := s.Configure(code); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/configure", http.StatusSeeOther)
}

func (s *GoogleDriveStorage) Configure(setup interface{}) error {
    // GoogleDriveStorage setup requires a string for setup
    code, ok := setup.(string)
    if !ok {
        return fmt.Errorf("Failed to setup Google Drive Client: Passed Code was not a string")
    }
    // Exchange code to get oauth token
    tok, err := s.Config.Exchange(context.Background(), code)
    if err != nil {
        return fmt.Errorf("Failed to create Google Drive Client: %v", err)
    }
    // Create Oatuh Client
    client := s.Config.Client(context.Background(), tok)
    // Use Client to create a google Drive service
    srv, err := drive.New(client)
    if err != nil {
        return fmt.Errorf("Failed to create Google Drive Client: %v", err)
    }
    s.Lock.Lock()
    defer s.Lock.Unlock()
    s.Service = srv
    return nil
}

func (s *GoogleDriveStorage) loadCredentials() (error){
    b, err := ioutil.ReadFile(s.ConfigPath)
    if err != nil {
        return err
    }

    config, err := google.ConfigFromJSON(b, drive.DriveScope)
    if err != nil {
        return err
    }

    s.Config = config
    return nil
}


func (s *GoogleDriveStorage) NewFolder(name string, parentId string) (string, error) {
    var parents []string
    if parentId != "" {
        parents = []string{parentId}
    }
    directory, err := s.Service.Files.Create(&drive.File{
        Name:name,
        MimeType: MIME_TYPE_GOOGLE_DRIVE_FOLDER,
        Parents: parents,
    }).Do()
    return directory.Id, err
}

func (s *GoogleDriveStorage) NewFile(name string, parentId string, mimeType string, content io.Reader) (string, error){
    var parents []string
    if parentId != "" {
        parents = []string{parentId}
    }
    file, err := s.Service.Files.Create(&drive.File{
        Name:name,
        MimeType: mimeType,
        Parents: parents,
    }).Media(content).Do()
    return file.Id, err
}

func (s *GoogleDriveStorage) Update(id string, content io.Reader) error {
    fmt.Println("Before: " , id)
    file, err := s.Service.Files.Update(id, nil).Media(content).Do()
    fmt.Println("After: ",file.Id)
    return err
}

func (s *GoogleDriveStorage) Delete(id string) error {
    return s.Service.Files.Delete(id).Do()
}

func (s *GoogleDriveStorage) Get(id string) (io.ReadCloser, error){
    res, err := s.Service.Files.Get(id).Download()
    return res.Body, err
}

func (s *GoogleDriveStorage) List(parentId string) ([]string, error){
    return s.Search(Query{ParentId: parentId})
}

func (s *GoogleDriveStorage) Search(query Query) ([]string, error){
    buildQuery := buildQuery(query)
    fmt.Printf("Built query %s from struct %v\n", buildQuery, query)
    r, err := s.Service.Files.List().PageSize(1000).
                Fields("nextPageToken, files(id, name)").
                Q(buildQuery).
                Do()
    if err != nil {
        return nil, err
    }
    res := make([]string,0, len(r.Files))
    for _, file := range r.Files {
        res = append(res, file.Id)
    }
    return res, nil
}

type Query struct {
    Name string
    MimeType string
    ParentId string
}

func buildQuery(query Query) string {
    queryString := ""
    if query.Name != ""{
        queryString = appendQuery(queryString, "name='" + query.Name + "'")
    }
    if query.MimeType != ""{
        queryString = appendQuery(queryString, "mimeType='" + query.MimeType + "'")
    }
    if query.ParentId != "" {
        queryString = appendQuery(queryString, "'" + query.ParentId + "' in parents")
    }
    return queryString
}

func appendQuery(query, value string) string{
    if query == "" {
        return value
    }
    return query + " and " + value
}
