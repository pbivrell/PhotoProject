package storage

import "fmt"

const (
    STORAGE_TYPE_LOCAL = "local"
    STORAGE_TYPE_GOOGLE_DRIVE = "googleDrive"
    STORAGE_MODE_WRITE = "w"
    STORAGE_MODE_READ = "r"
    STORAGE_CACHE_STRATAGEY_PERMANENT = "permanent"
    STORAGE_CACHE_STRATAGEY_MRU = "mru"
    STORAGE_CACHE_STRATAGEY_NEVER = "never"
    STORAGE_NAME_DEFAULT = "anonymous"
)

type FilePersmissions interface{}

type FileSearch struct{
    Name FileName
    Id FileId
    MimeType FileMimeType
}

type FileName interface{}

type FileId interface{}

type FileMimeType interface{
    IsFolder() bool
    Type() interface{}
}

type FileContents interface{}

type StorageConfiguration interface{}

type Storage interface {
    Create(FileSearch, FilePermissions, FileMimeType) (Id,  err)
    Configure(StorageConfiguration) string
    Configured() bool
    Delete(FileSearch) err
    Get(FileSearch) (FileContents, err)
    Name() string
    CacheMode() string
    ReadWriteMode() string
    Type() string
    Path() string
    List(Search) ([]Id, err)
    Update(Search, FileContents) err
}

func Create(search FileSearch, perms FilePermissions, mimeType FileMimeType) (Id, err){
    
}
