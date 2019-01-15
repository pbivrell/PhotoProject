package photoProject

import (
        "fmt"
        //"log"
        "net/http"
        //"os"
        "bufio"
        "bytes"
        "image/jpeg"
        "encoding/json"
        "google.golang.org/api/drive/v3"
    
        "github.com/disintegration/imageorient"
        "github.com/disintegration/imaging"

)

func CreatePageConfig(srv *drive.Service, projectID string, data LoadData, w http.ResponseWriter) string {
    b, err := json.Marshal(data)
    if err != nil {
        fmt.Fprintf(w, "Unable to create page configuration: %v\n", err)
        return ""
    }
    configFile, err := srv.Files.Create(&drive.File{Name: "config.json", MimeType:"application/json", Parents: []string{projectID}}).Media(bytes.NewReader(b)).Do()
    srv.Permissions.Create(configFile.Id, &drive.Permission{Role: "reader", Type: "anyone"}).Do()
    if err != nil {
        fmt.Fprintf(w, "Unable to create page configuration: %v\n", err)
        return ""
    }
    return configFile.Id
}

const ImageProjectRoot = "PhotoProjectRoot"
const ImageProjectIndex = "PhotoProjectIndex.json"

func CreateRootFile(srv *drive.Service) (string, string, error) {
    r, err := srv.Files.List().PageSize(1).
                Fields("nextPageToken, files(id, name)").
                Q("name='"+ImageProjectIndex+"'").Do()
    if err != nil {
        return "","", fmt.Errorf("Error looking for index file: %v\n", err)
    }
    if len(r.Files) == 0 {
        rootDir, err := srv.Files.Create(&drive.File{Name:ImageProjectRoot, MimeType:"application/vnd.google-apps.folder"}).Do()
        if err != nil {
            return "", "", fmt.Errorf("Error creating root directory: %v\n", err)
        }
        indexFile, err := srv.Files.Create(&drive.File{Name:ImageProjectIndex , MimeType:"application/json", Parents: []string{rootDir.Id}}).Media(bufio.NewReader(bytes.NewBufferString("{ items: [] }"))).Do()
        
        if err != nil {
            return "", "", fmt.Errorf("Error creating index.json: %v\n", err)
        }
        srv.Permissions.Create(indexFile.Id, &drive.Permission{Role: "reader", Type: "anyone"}).Do()
        return rootDir.Id, indexFile.Id, nil
    }
    // This assumes that the photoProjectIndex file is in the root directory. This all falls apart if it is not
    // I think this is acceptable because the person running the server controls where the images are written to
    // and they would only be shooting themselves in the foot if this was wrong.
    return r.Files[0].Parents[0], r.Files[0].Id, nil
}

func CreateNewProject(srv *drive.Service, rootId, projectId string) (string, *drive.File, *drive.File, error) {
    project, err := srv.Files.Create(&drive.File{Name:projectId, MimeType:"application/vnd.google-apps.folder", Parents: []string{rootId}}).Do()
    if err != nil {
        return "", nil, nil, fmt.Errorf("Error creating project directory: %v\n", err)
    }
    tiny, err := srv.Files.Create(&drive.File{Name:"tiny", MimeType:"application/vnd.google-apps.folder", Parents: []string{project.Id}}).Do()
    if err != nil {
        return "", nil, nil, fmt.Errorf("Error creating tiny directory: %v\n", err)
    }
    big, err := srv.Files.Create(&drive.File{Name:"big", MimeType:"application/vnd.google-apps.folder", Parents: []string{project.Id}}).Do()
    if err != nil {
        return "", nil, nil, fmt.Errorf("Error creating big directory: %v\n", err)
    }
    return project.Id, tiny, big, nil
}

func ProcessImages(srv *drive.Service, rootId, folderId string) (string, []string, []string, error) {
    projectId, tiny, big, err := CreateNewProject(srv, rootId, folderId)
    if err != nil {
        return projectId, nil, nil, err
    }
    r, err := srv.Files.List().PageSize(1000).
                Fields("nextPageToken, files(id, name)").
                Q("mimeType='image/jpeg' and '"+ folderId + "' in parents").
                Do()
    if err != nil {
        return projectId, nil, nil, fmt.Errorf("Error listing provided urls files: %v\n", err)
    }
    big_images := make([]string, 0, len(r.Files))
    tiny_images := make([]string, 0, len(r.Files))
    for _, f := range r.Files {
        id, err := ProcessImage(srv, f, tiny.Id, true)
        if err == nil {
            tiny_images = append(tiny_images, id)
        }
        id, err = ProcessImage(srv, f, big.Id, false)
        if err == nil {
            big_images = append(big_images, id)
        }
    }
    return projectId, tiny_images, big_images, nil
}

func ProcessImage(srv *drive.Service, f *drive.File, parentID string, shrink bool) (string, error) {
    res, err := srv.Files.Get(f.Id).Download()
    if err != nil {
        return "", err
    }
    img, _, err := imageorient.Decode(res.Body)
    if err != nil {
        return "", err
    }
    if shrink {
        img = imaging.Resize(img, 20, 0, imaging.Lanczos)
    }
    var b bytes.Buffer
    imageWriter := bufio.NewWriter(&b)
    err = jpeg.Encode(imageWriter, img, nil)
    if err != nil {
        return "", err
    }
    
    imageFile, err := srv.Files.Create(&drive.File{Name: f.Name, MimeType:"image/jpeg", Parents: []string{parentID}}).Media(bufio.NewReader(&b)).Do()
    if err != nil {
        return "", err
    }
    srv.Permissions.Create(imageFile.Id, &drive.Permission{Role: "reader", Type: "anyone"}).Do()
    return imageFile.Id, nil
}
