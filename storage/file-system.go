package storage

import (
)

func main(){

}

type FileSystemStorage struct {

}

func newStorageFileSystem(mode,cacheStrat,name, path string) *GoogleDriveStorag {

}

func (storage *FileSystemStorage) Create(FileSearch, FilePermissions, FileMimeType) (Id,  err){

}

func (storage *FileSystemStorage) Configure(StorageConfiguration) err  {

}

func (storage *FileSystemStorage) Configured() bool {

}

func (storage *FileSystemStorage) Delete(FileSearch) err {

}

func (storage *FileSystemStorage) Get(FileSearch) (FileContents, err) {

}

func (storage *FileSystemStorage) Name() string {

}

func (storage *FileSystemStorage) CacheMode() string {

}

func (storage *FileSystemStorage) ReadWriteMode() string {

}

func (storage *FileSystemStorage) Type() string {

}

func (storage *FileSystemStorage) Path() string {

}

func (storage *FileSystemStorage) List(Search) ([]Id, err) {

}

func (storage *FileSystemStorage) Update(Search) err {

}
