package main

import (
    "io"
)

type Storage interface {
    // Storage Interaction Method
    NewFolder(name string, parentId ...string) (File, error)
    NewFile(name string, content io.Reader, parentId ...string) (File, error)
    IsFolder(id string) (bool, error)
    Update(id string, content io.Reader) (File, error)
    Delete(id string) error
    Get(id string) (io.Reader, error)
    List(parentId string) ([]File, error)
    Search(query File) ([]File, error)
}

type File struct {
    Name string
    ParentIds []string
    Id string
}
