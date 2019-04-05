package main

import (
    ""
)

const (
)

func NewMockStorage() *FileSystemWrapper {
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
