package storage

import (
)

func main(){

}

type GoogleDriveStorage struct {

}

func newStorageGoogleDrive(mode,cacheStrat,name, path string) *GoogleDriveStorag {

}

func (storage *GoogleDriveStorage) Create(FileSearch, FilePermissions, FileMimeType) (Id,  err){

}

func (storage *GoogleDriveStorage) Configure(StorageConfiguration) err  {

}

func (storage *GoogleDriveStorage) Configured() bool {

}

func (storage *GoogleDriveStorage) Delete(FileSearch) err {

}

func (storage *GoogleDriveStorage) Get(FileSearch) (FileContents, err) {

}

func (storage *GoogleDriveStorage) Name() string {

}

func (storage *GoogleDriveStorage) CacheMode() string {

}

func (storage *GoogleDriveStorage) ReadWriteMode() string {

}

func (storage *GoogleDriveStorage) Type() string {

}

func (storage *GoogleDriveStorage) Path() string {

}

func (storage *GoogleDriveStorage) List(Search) ([]Id, err) {

}

func (storage *GoogleDriveStorage) Update(Search) err {

}
