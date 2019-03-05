package main

type StorageStructure struct {
    Node StorageMetadata
    Children Childs
}

type Childs []*StorageStructure

type StorageMetadata struct {
    Id string
    Name string
    ParentId string
    MimeType string
}

func (s Childs) Len() int {
    return len(s)
}

func (s Childs) Less(i, j int) bool {
    return s[i].Node.Name < s[j].Node.Name
}

func (s Childs) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func NewStorageStructure() *StorageStructure {
    return &StorageStructure{Node: &StorageMetadata{}, Children: make([]*StorageStructure,0)}
}

func (s *StorageStructure) Search(search string) *StorageStructure {
    return s.depthLimtedSearch(search, -1)
}

func (s *StorageStructure) List(search string) *StorageStructure {
    return s.depthLimtedSearch(search, -1)
}

func (s *StorageStructure) depthLimitedSearch(search string, depth int) *StorageStructure {
    unexplored := []*StorageStructure{ s }
    remainElements := 1
    for len(unexplored) > 0 && depth == 0 {
        if unexplored[0].Node.Name == search {
            return unexplored[0]
        }else{
            unexplored := append(unexplored, unexplored[0].Children...)
        }
        unexplored[0] = nil
        unexplored = unexplored[1:]
        remainElements--
        if remainElements == 0 {
            depth--
            remainingElements = len(unexplored)
        }
    }
    return nil
}

func (s *StorageStructure) InsertByParentId(name, id, mimeType, parentId string, directory bool){

}

func (s *StorageStructure) InsertByPath(name, id, mimeType, path string, directory bool){
    
}
