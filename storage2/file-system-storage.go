package storage


import (
    "github.com/gorilla/mux"

    "fmt"
    "net/http"
    "sync"
    "io/ioutil"
    "io"
)

const (
    STORAGE_TYPE_FILE_SYSTEM= "local"
)

type Storage interface {
    ConfigureEndpoints(r *mux.Router)
    Configure() string
    Configured() bool
    NewFolder(name string, parentId string) (string, error)
    NewFile(name string, parentId string, ctype string, content io.Reader) (string, error)
    Update(id string, content io.Reader) error
    Delete(id string) error
    Get(id string) (io.ReadCloser, error)
    List(parentId string) ([]string, error)
    //Search(name string, options string...) ids []string
    //Exists(id string) bool
}

type FileSystemStorage struct {
    Root string
    VolumeId string
}

func NewFileSystemStorage(volumeId string, rootMountPath string) *FileSystemStorage {
    return &FileSystemStorage{
        Root: rootMountPath,
        VolumeId: volumeId
    }
}

func (s *FileSystemStorage) Configure() string {
    return fmt.Sprintf("Storage: File System \"%s\": Configured")
}

func (s *FileSystemStorage) Configured() bool {
    return true
}

func (s *FileSystemStorage) ConfigureEndpoints(r *mux.Router) {}

func (s *FileSystemStorage) NewFolder(name string, parentId string) (string, error) {
    dir := s.Root + parentId + "/"
    err := os.Mkdir(dir + name, 0x600)
    return name, err
}

func (s *FileSystemStorage) NewFile(name string, parentId string, mimeType string, content io.Reader) (string, error){
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

func (s *FileSystemStorage) Update(id string, content io.Reader) error {
    fmt.Println("Before: " , id)
    file, err := s.Service.Files.Update(id, nil).Media(content).Do()
    fmt.Println("After: ",file.Id)
    return err
}

func (s *FileSystemStorage) Delete(id string) error {
    return s.Service.Files.Delete(id).Do()
}

func (s *FileSystemStorage) Get(id string) (io.ReadCloser, error){
    res, err := s.Service.Files.Get(id).Download()
    return res.Body, err
}

func (s *FileSystemStorage) List(parentId string) ([]string, error){
    r, err := s.Service.Files.List().PageSize(1000).
                Fields("nextPageToken, files(id, name)").
                Q(parentId + "' in parents").
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
