package main

import (
    "os"
    "fmt"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "sync"
    "math/rand"
    "time"

    "github.com/gorilla/mux"
    "github.com/pbivrell/ManagedMap"

    "golang.org/x/crypto/bcrypt"
)

const (
    // These constants define the mapping between indeces of the WriteStorage and
    // there meaning. The server can be configured have 3 difference storage volumes based
    // on likely traffic patterns. The root page will likely have the heaviest traffic so
    // it can be configured to be served by a quick-to-access storage volume. Similarily with
    // the project json pages. If a level has not been setup the storage is required to be delegated
    // to the level immediately below it ie if ROOT_JSON was not setup then the ROOT_JSON data will 
    // be stored in the PROJECT_JSONS if it exists.
    WRITE_STORAGE_PHOTOS = 0
    WRITE_STORAGE_PROJECT_JSONS = 1
    WRITE_STORAGE_ROOT_JSON = 2
    // Variables for loading the configuration file
    CONFIG_ENV = "PHOTO_PROJECT_CONFIG"
    CONFIG_DEFAULT_PATH = "./config.json"
    // Token Size
    TOKEN_SIZE = 32
)

type configuration struct {
    IP     string `json:"ip,omitempty"`
    Port  string `json:"port,omitempty"`
    AdmitPassword string `json:"admitPassword, omitempty"`
    StaticDir string `json:"staticDir"`
    WriteStorage []struct {
        Name  string `json:"name"`
        Type  string `json:"type"`
        Path  string `json:"path,omitempty"`
        Cache string `json:"cache,omitempty"`
        Extra string `json:"extra"`
    } `json:"writeStorage"`
}

func LoadConfig() (*configuration, error) {
    configPath := CONFIG_DEFAULT_PATH
    if path := os.Getenv(CONFIG_ENV); path != "" {
        configPath = path
    }
    configFile, err := os.Open(configPath)
    if err != nil {
        return nil, fmt.Errorf("Unable to open config file: %v\n",err)
    }
    defer configFile.Close()
    data, err := ioutil.ReadAll(configFile)
    if err != nil {
        return nil, fmt.Errorf("Unable to read config file contents: %v\n", err)
    }
    var config configuration
    err = json.Unmarshal(data, &config)
    if err != nil {
        return nil, fmt.Errorf("Unable to process config file contents: %v\n", err)
    }
    if config.IP == "" {
        config.IP = "localhost"
    }
    if config.Port == "" {
        config.Port = "8080"
    }
    return &config, nil
}

type Server struct {
    AdmitHash string
    WriteStorage []Storage
    StaticDir string
    Lock *sync.RWMutex
    TokenMap *ManagedMap.ManagedMap
    IP string
    Port string
}

func (s *Server) Configured() bool {
    configured := true
    for _, storage := range s.WriteStorage {
        configured = configured && storage.Configured()
    }
    s.Lock.RLock()
    defer s.Lock.RUnlock()
    return configured
}

func NewServer(config *configuration, hash string, r *mux.Router) *Server {
    writeStorage := make([]Storage, 0)
    for _, v := range config.WriteStorage {
        storage := NewStorage(v.Type, v.Name, v.Extra)
        storage.ConfigureEndpoints(r)
        writeStorage = append(writeStorage, storage)
    }
    return &Server{
        AdmitHash: hash,
        WriteStorage: writeStorage,
        StaticDir: config.StaticDir,
        Lock: &sync.RWMutex{},
        TokenMap: ManagedMap.NewManagedMap(),
        IP: config.IP,
        Port: config.Port,
    }
}

var server *Server

// Utility Functions
func HashPassword(password string) (string, error) {
    if password == "" {
        return "", nil
    }
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func init() {
    rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GetRandomToken() string {
    b := make([]rune, TOKEN_SIZE)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

func main(){
    config, _ := LoadConfig()
    hash, _ := HashPassword(config.AdmitPassword)
    r := mux.NewRouter()
    server = NewServer(config, hash, r)
    r = configureRouter(r)
    fmt.Println(http.ListenAndServe(server.IP+":"+ server.Port, r))
}
