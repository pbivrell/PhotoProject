package main

import (
    "google.golang.org/api/drive/v3"
    "golang.org/x/net/context"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"

    "net/http"
    "io/ioutil"
)

type GoogleDriveOauth struct {
    GenerateState func() string
    ValidateState func(string) bool
    CallbackUrlPath string
    CallbackSuccessHandler func(http.ResponseWriter, *http.Request)
    CallbackFailureHandler func(http.ResponseWriter, *http.Request)
    GetCredentials func() ([]byte, error)
    Scopes []string
    AccessType oauth2.AuthCodeOption
    config *oauth2.Config
    client *http.Client
    clientSetup chan struct{}
}

func NewGoogleDriveOauth() *GoogleDriveOauth {
    return &GoogleDriveOauth{
        GenerateState: func() string {
            return "state"
        },
        ValidateState: func(state string) bool {
            return state == "state"
        },
        CallbackUrlPath: "/configureGDS-test",
        GetCredentials: func ()([]byte, error) {
            return ioutil.ReadFile("./credentials.json")
        },
        CallbackSuccessHandler: func(w http.ResponseWriter, r *http.Request) {
            http.Redirect(w,r,"/", http.StatusSeeOther)
        },
        CallbackFailureHandler: func(w http.ResponseWriter, r *http.Request) {
            http.Error(w, "Failed to configure Google Drive Storage using Oauth2", http.StatusInternalServerError)
        },
        Scopes: []string{drive.DriveScope},
        AccessType: oauth2.AccessTypeOffline,
        config: nil,
        client: nil,
        clientSetup: make(chan struct{}),
    }
}

type Endpoint struct {
    Path string
    Handler func(http.ResponseWriter, *http.Request)
}

func (s *GoogleDriveOauth) GetCallbackEndpoint() Endpoint {
    return Endpoint{ Path: s.CallbackUrlPath, Handler: s.CallbackHandler}
}

func (s *GoogleDriveOauth) CallbackHandler(w http.ResponseWriter, r *http.Request){
    state := r.FormValue("state")
    if !s.ValidateState(state) {
        s.CallbackFailureHandler(w,r)
        return
    }
    code := r.FormValue("code")
    if err := s.CreateClient(code); err != nil {
        s.CallbackFailureHandler(w,r)
        return
    }
    s.clientSetup <- struct{}{}
    s.CallbackSuccessHandler(w,r)
}

func (s *GoogleDriveOauth) GetAuthUrl() string {
    return s.config.AuthCodeURL(s.GenerateState(), s.AccessType)
}

func (s *GoogleDriveOauth) LoadCredentialsFromJSON() (error){
    bytes, err := s.GetCredentials()
    if err != nil {
        return err
    }
    config, err := google.ConfigFromJSON(bytes, s.Scopes...)
    if err != nil {
        return err
    }
    s.config = config
    return nil
}

func (s *GoogleDriveOauth) CreateClient(code string) error{
    // Exchange code to get oauth token
    tok, err := s.config.Exchange(context.Background(), code)
    if err != nil {
        return err
    }
    s.client = s.config.Client(context.Background(), tok)
    return nil
}

func (s *GoogleDriveOauth) GetClient() *http.Client {
    <-s.clientSetup
    return s.client
}

func DoOauth() (Endpoint, string, <-chan *http.Client) {
    oauth := NewGoogleDriveOauth()
    oauth.LoadCredentialsFromJSON()
    client := make(chan *http.Client)
    go func() {
        client <- oauth.GetClient()
    }()
    return oauth.GetCallbackEndpoint(), oauth.GetAuthUrl(), client
}

/* Example Usage
func main(){
    oauth := NewGoogleDriveOauth()
    oauth.LoadCredentialsFromJSON()
    endpoint := oauth.GetCallbackEndpoint()
    fmt.Println(oauth.GetAuthUrl())
    http.HandleFunc(endpoint.Path, endpoint.Handler)
    go http.ListenAndServe(":8080", nil)
    oauth.GetClient()
    fmt.Println("Client Setup")
} */
