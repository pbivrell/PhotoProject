package main

import (
    "io"
    "fmt"
    "strings"
)

type FileMetadata struct {
    Name string
    Id string
    IsDir bool
    Parent *FileMetadata
    Children []*FileMetadata
}

func (f *FileSystemWrapper) PathSearch(path string) ([]File, error){
    prev := []File{}
    for i,v := range strings.Split(path, "/") {
        fs, err := f.s.Search(File{Name: v})
        if err != nil {
            return nil, err
        }
        if i == 0 {
            prev = fs
            continue 
        }
        newFs := make([]File, 0)
        if newFs = In(prev,fs); len(newFs) == 0 {
            return []File{}, nil
        }
        prev = newFs
    }
    return prev, nil
}

func In(prev []File, curr []File) []File {
    res := make([]File, 0)
    if len(prev) == 0 || len(curr) == 0 {
        return res
    }

    for _,v := range prev {
        for _, v2 := range curr {
            for _, v3 := range v2.ParentIds {
                if v.Id == v3 {
                    res = append(res, v2)
                }
            }
        }
    }
    return res
}

type FileSystemWrapper struct {
    s Storage
    root *FileMetadata
    ids map[string]*FileMetadata
}

func NewFileSystemWrapper(originalStorage Storage, mountId string) (*FileSystemWrapper, error) {
    ids := make(map[string]*FileMetadata)
    isRoot, err := originalStorage.IsRoot(mountId)
    if err != nil {
        return nil, UnderlayingStorageError(err)
    }
    if isRoot {
        fs := &FileSystemWrapper{
            s: originalStorage,
            ids: ids,
            root: &FileMetadata{
                Name: "",
                Id: mountId,
                IsDir: true,
                Parent: nil,
                Children: make([]*FileMetadata,0),
            },
        }
        fs.ids[mountId] = fs.root
        return fs, nil
    }
    rootFile, err := originalStorage.GetMetadata(mountId)
    if err != nil {
        return nil, UnderlayingStorageError(err)
    }
    isDir, err := originalStorage.IsFolder(mountId)
    if err != nil {
        return nil, UnderlayingStorageError(err)
    }
    if !isDir {
        return nil, fmt.Errorf("Failed to mount at id [%s] original storage is not a directory.")
    }
    fs := &FileSystemWrapper{
        s: originalStorage,
        ids: ids,
        root: &FileMetadata{
            Name: "",
            Id: rootFile.Id,
            IsDir: isDir,
            Parent: nil,
            Children: make([]*FileMetadata, 0),
        },
    }
    fs.ids[fs.root.Id] = fs.root
    return fs, nil
}

func UnderlayingStorageError(err error) error {
    return fmt.Errorf("Underlaying storage error: %s", err.Error());
}

func (f *FileSystemWrapper) Crawl() error {
    dirs :=[]*FileMetadata{ f.root }
    for len(dirs) > 0 {
        curr := dirs[0]
        dirs = dirs[1:]
        //fmt.Println("Crawling")
        //fmt.Println(curr.Name)
        //fmt.Println(curr.Id)
        children,err := f.s.List(curr.Id)
        //fmt.Println(len(children))
        if err != nil {
            return UnderlayingStorageError(err)
        }
        for _,v := range children {
            isDir, err := f.s.IsFolder(v.Id)
            if err != nil {
                return UnderlayingStorageError(err)
            }
            child := &FileMetadata{
                Name: v.Name,
                Id: v.Id,
                IsDir: isDir,
                Parent: curr,
                Children: make([]*FileMetadata, 0),
            }
            curr.Children = append(curr.Children, child)
            f.ids[child.Id] = child
            if isDir {
                dirs = append(dirs, child)
            }
        }
        //fmt.Println(curr.Children)
    }
    return nil
}

func search(c []*FileMetadata, name string) *FileMetadata {
    for _,p := range c {
        //fmt.Println("'"+p.Name+"' is '"+ name +"'")
        //fmt.Println(p.Name == name)
        if p.Name == name {
            //fmt.Println("Returning")
            return p
        }

    }
    return nil
}

func (f *FileSystemWrapper) PathToId(path string) []string {
    if i := strings.LastIndex(path, "/"); i != -1 && len(path) -1 == i{
        //fmt.Print("len %d, index %d\n", len(path), strings.LastIndex(path, "/"))
        path = path[:len(path)-1]
    }
    if strings.Index(path, "/") == 0 {
        path = path[1:]
    }
    //fmt.Println("Path:", path)
    //fmt.Println("root",f.root.Name)
    //fmt.Println(path == f.root.Name)
    if path == f.root.Name {
        return []string{f.root.Id}
    }
    //fmt.Println("Here")
    child := f.root
    for _, v := range strings.Split(path, "/") {
        //fmt.Println("c",child)
        //fmt.Println("v",v)
        //if !child.IsDir && child.Name == v {
        //    return []string{ child.Id }
        //}
        child = search(child.Children, v)
        //fmt.Println(child)
        if child == nil {
            return []string{}
        }
    }
    //fmt.Println(child)
    return []string{ child.Id }
}

//func (s *GoogleDriveStorage) NewFolder(name string, parentIds ...string) (File, error) {
func (f *FileSystemWrapper)  NewFolder(name string, parentIds ...string) (File, error){
    return File{}, nil
}

func (f *FileSystemWrapper) IsFolder(id string) (bool, error) {
    file, has := f.ids[id]
    if !has {
        return false, fmt.Errorf("Couldn't find file with id %s", id)
    }
    return file.IsDir, nil
}

func (f *FileSystemWrapper) NewFile(name string, content io.Reader, parentId ...string) (File, error){
    return File{}, nil
}

func (f *FileSystemWrapper) Update(id string, content io.Reader) (File, error){
    return File{}, nil
}

func (f *FileSystemWrapper) Delete(id string) error {
    return nil
}

func (f *FileSystemWrapper) Get(id string) (io.Reader, error) {
    return f.s.Get(id)
}
func (f *FileSystemWrapper) GetFSMetadata(id string) (*FileMetadata, error) {
    file, has := f.ids[id]
    if !has {
        return &FileMetadata{}, fmt.Errorf("Couldn't find file with id %s", id)
    }
    return file, nil
}

func (f *FileSystemWrapper) GetMetadata(id string) (File, error) {
    file, has := f.ids[id]
    if !has {
        return File{}, fmt.Errorf("Couldn't find file with id %s", id)
    }
    return File{Name: file.Name, Id: file.Id, ParentIds:[]string{ }}, nil
}

func (f *FileSystemWrapper) IsRoot(id string) (bool, error){
    return false, nil
}

func (f *FileSystemWrapper) List(parentId string) ([]File, error){
    file, has := f.ids[parentId]
    if !has {
        return []File{}, fmt.Errorf("Couldn't find file with id %s", parentId)
    }
    children := make([]File, len(file.Children))
    for i, v := range file.Children {
        children[i] = File{Name: v.Name, Id: v.Id, ParentIds:[]string{ parentId, }}
    }
    return children, nil
}

func (f *FileSystemWrapper) Search(query File) ([]File, error){
    return []File{}, nil
}
