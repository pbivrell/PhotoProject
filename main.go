package main

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "os"
        "image/jpeg"

        "golang.org/x/net/context"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/drive/v3"
    
        "github.com/disintegration/imageorient"
        "github.com/disintegration/imaging"

)

const (
    tokenFile = "token.json"
    credentialsFile = "credentials.json"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
        // The file token.json stores the user's access and refresh tokens, and is
        // created automatically when the authorization flow completes for the first
        // time.
        tok, err := tokenFromFile(tokenFile)
        if err != nil {
        
        }
        return config.Client(context.Background(), tok)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
        f, err := os.Open(file)
        if err != nil {
                return nil, err
        }
        defer f.Close()
        tok := &oauth2.Token{}
        err = json.NewDecoder(f).Decode(tok)
        return tok, err
}

func GetDriveFolder(folderId string, uuid string, srv *drive.Service) error{
        r, err := srv.Files.List().PageSize(1000).
                Fields("nextPageToken, files(id, name)").
                Q("mimeType='image/jpeg' and '"+ folderId + "' in parents").
                Do()
        
        if err != nil {
                log.Fatalf("Unable to retrieve files: %v", err)
                return err
        }
        
        fmt.Printf("FOUND: %d\n", len(r.Files))
        
        srv.Files.Create


        for _, i := range r.Files {
            res, err := srv.Files.Get(i.Id).Download()
            if err != nil {
                log.Fatalf("Error: %v", err)
                return err
            }
            img, _, err := imageorient.Decode(res.Body)
            if err != nil {
                log.Fatalf("imageorient.Decode failed: %v", err)
            }

            f, err := os.Create("test/"+i.Name)
            if err != nil {
                log.Fatalf("os.Create failed: %v", err)
            }
            err = jpeg.Encode(f, img, nil)
            if err != nil {
                log.Fatalf("jpeg.Encode failed: %v", err)
            }
            img = imaging.Resize(img, 20, 0, imaging.Lanczos)
            f, err = os.Create("test/tiny_"+i.Name)
            if err != nil {
                log.Fatalf("os.Create failed: %v", err)
            }
            err = jpeg.Encode(f, img, nil)
            if err != nil {
                log.Fatalf("jpeg.Encode failed: %v", err)
            }

            
        }

        return nil
}

func Authenticate() (*drive.Service, error) {
        b, err := ioutil.ReadFile(credentialsFile)
        if err != nil {
                log.Fatalf("Unable to read client secret file: %v", err)
                return nil, err
        }

        // If modifying these scopes, delete your previously saved token.json.
        config, err := google.ConfigFromJSON(b, drive.DriveReadonlyScope)
        if err != nil {
                log.Fatalf("Unable to parse client secret file to config: %v", err)
                return nil, err
        }
        client := getClient(config)

        srv, err := drive.New(client)
        if err != nil {
                log.Fatalf("Unable to retrieve Drive client: %v", err)
                return nil, err
        }
        return srv, nil

}

func main() {
        srv, err := Authenticate()
        if err != nil || srv == nil  {
            fmt.Println("BADNESS")
            return
        }
        parentID := "14RcDryyKPvN8pkkyI3fZPzCUcZiGS34q" 
       //parentID := "1l-cQFuwP_3mrsVSa-rFGxWDhTv8KH4qa"
        //parentID := "1CQ3JdCxFNj-MYOvxYeiGVqYjZNF7gelR"
        err = GetDriveFolder(parentID, "test2", srv)
    
}

