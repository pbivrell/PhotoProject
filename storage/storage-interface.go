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
    Update(Search) err
}

func NewStorage(type, mode, cacheStrat, name, path string) (*Storage, err) {
    if cacheStrat != STORAGE_CACHE_STRATAGEY_PERMANENT || cache != STORAGE_CACHE_STRATAGEY_MRU {
        cacheStrat = STROAGE_CACHE_STRATAGEY_NEVER
    }
    if name == "" {
        name == STORAGE_NAME_DEFAULT
    }
    if type == STORAGE_TYPE_LOCAL{
        return newStorageLocal(mode, cacheStrat, name, path)
    }else if type == STORAGE_TYPE_GOOGLE_DRIVE {
        return newStorageGoogleDrive(mode,cacheStrat,name, path)
    }
    return nil, fmt.Errorf("Couldn't create storage of type %s. Does not match any supported storage types\n", type)
}
