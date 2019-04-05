package main

import (
    ""
)

const (
    MAX_RETRY = 3
)

type FileMetadata struct {
    Name string
    Id string
    IsDir bool
    Parents map[string]*FileMetadata
    Children map[string]*FileMetadata
}

type FileSystemWrapper struct {
    storage *Storage
    mountPath string
    root *FileMetadata
}


func NewFileSystemWrapper(originalStorage *Storage, mountPath string) *FileSystemWrapper {
    return &{storage: originalStorage, mountPath: mouthPath,}
}

func (f *FileSystemWrapper) CrawlSource() {
}

func (f *FileSystemWrapper) NewFolder(name string, parentId ...string) (File, error){
    
}

func (f *FileSystemWrapper) IsFolder(id string) (bool, error) {

}

func (f *FileSystemWrapper) NewFile(name string, content io.Reader, parentId ...string) (File, error){

}

func (f *FileSystemWrapper) Update(id string, content io.Reader) (File, error){

}

func (f *FileSystemWrapper) Delete(id string) error {

}

func (f *FileSystemWrapper) Get(id string) (io.Reader, error) {

}

func (f *FileSystemWrapper) List(parentId string) ([]File, error){

}

func (f *FileSystemWrapper) Search(query File) ([]File, error){

}
