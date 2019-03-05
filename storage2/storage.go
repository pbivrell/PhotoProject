package main


import (
    "github.com/gorilla/mux"
)

const (
    MIME_TYPE_JSON = "application/json"
    MIME_TYPE_JPEG= "image/jpeg"
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
