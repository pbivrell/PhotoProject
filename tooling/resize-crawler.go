package main

import (
    "net/http"
    "fmt"
    "github.com/disintegration/imaging"
    "github.com/disintegration/imageorient"
    "bytes"
    "image/jpeg"
)

func main(){
    _, endpoint, authUrl, client := DoOauth()
    fmt.Println(authUrl)
    http.HandleFunc(endpoint.Path, endpoint.Handler)
    go http.ListenAndServe(":8080", nil)
    s, err := NewGoogleDriveStorage(<-client)
    //res, err := s.PathSearch("/Photos")
    //if err != nil {
    //    fmt.Println("Couldn't find \"/Photos\"", err)
    //    return
    //}
    small, err := s.GetMetadata("1C6vRZYTZh8tFqPunBtPRlT0f96MzaUE4")
    if err != nil {
        fmt.Println("Could not find Photos-Small:", err)
        return
    }
    normal, err := s.GetMetadata("1l_V63rO9ZkzluB1wIyISK8Yfbp6hBFjr")
    if err != nil {
        fmt.Println("Could not find Photos-Normal:", err)
        return
    }
    crawl(s, small, normal, 20)
}

func crawl(s Storage, to File, from File, size int) {
    dirs := []struct{ F File; P string } {{from, to.Id} }
    for len(dirs) > 0 {
        parent := dirs[0]
        dirs = dirs[1:]
        fmt.Println("Proccessing: ", parent.F.Name)
        newParent, err := s.NewFolder(parent.F.Name, parent.P)

        if err != nil {
            fmt.Println("Error creating resize parent:", parent.F.Name, ":", err)
            return
        }
        children, err := s.List(parent.F.Id)
        if err != nil {
            fmt.Println("Error listing children:", err)
            return
        }
        fmt.Printf("Children: %d:", len(children))
        for i, v := range children {
            if i % 5 == 0 {
                fmt.Printf("*")
            }
            if folder, err := s.IsFolder(v.Id); err == nil && folder {
                dirs = append(dirs, struct{F File; P string}{F:v, P:newParent.Id})
            }else {
                resize(s, v, newParent, size)
            }
        }
        fmt.Println("")
    }
}

func resize(s Storage, cp File, parent File, height int) {
    data, err := s.Get(cp.Id)
    img, _, err := imageorient.Decode(data)
    if err != nil {
        fmt.Printf("Err %v\n", err)
        return
    }
    img = imaging.Resize(img, height, 0, imaging.Lanczos)
    buffer := new(bytes.Buffer)
    if err := jpeg.Encode(buffer, img, nil); err != nil {
        fmt.Printf("unable to encode image:%v\n", err)
        return
    }
    _, err = s.NewFile(cp.Name, buffer, parent.Id)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
