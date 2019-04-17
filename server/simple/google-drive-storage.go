package main

import (
    "io"
    "net/http"
    "google.golang.org/api/drive/v3"
)


const (
    MIME_TYPE_GOOGLE_DRIVE_FOLDER = "application/vnd.google-apps.folder"
)

type GoogleDriveStorage struct {
    service *drive.Service
    Public bool
}

func NewGoogleDriveStorage(oauthClient *http.Client) (*GoogleDriveStorage, error) {
    srv, err := drive.New(oauthClient)
    if err != nil {
        return nil, err
    }
    return &GoogleDriveStorage{ service: srv, Public: true}, nil
}

func (s *GoogleDriveStorage) NewFolder(name string, parentIds ...string) (File, error) {
    directory, err := s.service.Files.Create(&drive.File{
        Name:name,
        MimeType: MIME_TYPE_GOOGLE_DRIVE_FOLDER,
        Parents: parentIds,
    }).Do()
    return File{Name:name, Id: directory.Id, ParentIds: parentIds}, err
}

func (s *GoogleDriveStorage) IsFolder(id string) (bool, error) {
    file, err := s.GetGoogleMetadata(id)
    return file.MimeType == MIME_TYPE_GOOGLE_DRIVE_FOLDER, err
}

func (s *GoogleDriveStorage) IsRoot(id string) (bool, error) {
    return id == "root", nil
}

type GoogleDriveMetadata struct {
    Name string
    Id string
    ParentIds []string
    MimeType string
}

func (s *GoogleDriveStorage) GetGoogleMetadata(id string) (GoogleDriveMetadata, error) {
    file, err := s.service.Files.Get(id).Do()
    if err != nil {
        return GoogleDriveMetadata{}, err
    }
    return GoogleDriveMetadata{Name: file.Name, Id: file.Id, ParentIds: file.Parents, MimeType: file.MimeType}, err
}

func (s *GoogleDriveStorage) GetMetadata(id string) (File, error) {
    metadata, err := s.GetGoogleMetadata(id)
    return File{Name: metadata.Name, Id: metadata.Id, ParentIds: metadata.ParentIds}, err
}

func (s *GoogleDriveStorage) NewFile(name string, content io.Reader, parentIds ...string) (File, error){
    file, err := s.service.Files.Create(&drive.File{
        Name:name,
        Parents: parentIds,
    }).Media(content).Do()
    if err != nil {
        return File{}, err
    }
    return File{Name:name, Id:file.Id, ParentIds: parentIds} , err
}

func (s *GoogleDriveStorage) Update(id string, content io.Reader) (File, error) {
    file, err := s.service.Files.Update(id, nil).Media(content).Do()
    return File{Name: file.Name, Id: file.Id, ParentIds: file.Parents}, err
}

func (s *GoogleDriveStorage) Delete(id string) error {
    return s.service.Files.Delete(id).Do()
}

func (s *GoogleDriveStorage) Get(id string) (io.Reader, error){
    var res *http.Response
    var err error
    if s.Public {
        //fmt.Println("Starting...")
        res, err = http.Get("https://drive.google.com/uc?export=view&id=" + id)
        //fmt.Println("Done")
    }else{
        res, err = s.service.Files.Get(id).Download()
    }
    return res.Body, err
}

func (s *GoogleDriveStorage) List(parentId string) ([]File, error){
    return s.Search(File{ParentIds: []string{parentId}})
}

func (s *GoogleDriveStorage) Search(data File) ([]File, error){
    query := buildQuery(data)
    r, err := s.service.Files.List().PageSize(1000).
    Fields("nextPageToken, files(id, name, parents)").
    Q(query).
    Do()
    if err != nil {
        return nil, err
    }
    res := make([]File,0)
    for _, file := range r.Files {
        res = append(res, File{Name: file.Name, Id: file.Id, ParentIds: file.Parents})
    }
    return res, nil
}


func buildQuery(data File) string {
    query := "trashed = false"
    for _, v := range data.ParentIds {
        if v != "" {
            query = appendAnd(query, "'" + v + "' in parents")
        }
    }
    if data.Name != "" {
        query = appendAnd(query, "name='" + data.Name + "'")
    }
    return query
}

func appendAnd(query string, newValue string) string {
    if query == "" {
        return newValue
    }
    return query + " and " + newValue
}

/*func main(){
    // Setup Oauth Stuff
    //oauth := NewGoogleDriveOauth()
    //oauth.LoadCredentialsFromJSON()
    //endpoint := oauth.GetCallbackEndpoint()
    //fmt.Println(oauth.GetAuthUrl())
    endpoint, authUrl, client := DoOauth()
    fmt.Println(authUrl)
    http.HandleFunc(endpoint.Path, endpoint.Handler)
    go http.ListenAndServe(":8080", nil)
    storage, err := NewGoogleDriveStorage(<-client)
    fmt.Println("Client Setup")
    if err != nil {
        fmt.Printf("%v\n", err)
    }
    l, _ := storage.PathSearch("a/b/c")
    for _, v:= range l {
        fmt.Println(v)
    }
    fmt.Println("Not path")
    l, _ = storage.Search(File{Name:"C"})
    for _, v:= range l {
        fmt.Println(v)
    }

    top,_ := storage.NewFolder("TEST FOLDER")
    fmt.Println(top)
    next,_ := storage.NewFolder("APPLE", top.Id)
    fmt.Println(next)
    n3,_ := storage.NewFile("pear.txt", strings.NewReader("apple"), top.Id)
    n2,_ := storage.NewFile("lime.txt", strings.NewReader("pear"), next.Id)
    storage.NewFile("b.json", strings.NewReader("{\"P\":\"Zebra\"}"), next.Id,top.Id)
    storage.Update(n2.Id, strings.NewReader("lemon"))
    storage.Delete(n3.Id)
    //content, _ := storage.Get(n2.Id)
    //buf := new(bytes.Buffer)
    //buf.ReadFrom(content)
    //fmt.Println(buf.String())
    l,_ := storage.Search(File{Name: "APPLE"})
    for _, v := range l {
        fmt.Println(v)
    }
}*/
