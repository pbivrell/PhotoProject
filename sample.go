package main

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "os"

        "golang.org/x/net/context"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/drive/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
        // The file token.json stores the user's access and refresh tokens, and is
        // created automatically when the authorization flow completes for the first
        // time.
        tokFile := "token.json"
        tok, err := tokenFromFile(tokFile)
        if err != nil {
                tok = getTokenFromWeb(config)
                saveToken(tokFile, tok)
        }
        return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
        authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
        fmt.Printf("Go to the following link in your browser then type the "+
                "authorization code: \n%v\n", authURL)

        var authCode string
        if _, err := fmt.Scan(&authCode); err != nil {
                log.Fatalf("Unable to read authorization code %v", err)
        }

        tok, err := config.Exchange(context.TODO(), authCode)
        if err != nil {
                log.Fatalf("Unable to retrieve token from web %v", err)
        }
        return tok
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

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
        fmt.Printf("Saving credential file to: %s\n", path)
        f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
        if err != nil {
                log.Fatalf("Unable to cache oauth token: %v", err)
        }
        defer f.Close()
        json.NewEncoder(f).Encode(token)
}

func main() {
        b, err := ioutil.ReadFile("credentials.json")
        if err != nil {
                log.Fatalf("Unable to read client secret file: %v", err)
        }

        // If modifying these scopes, delete your previously saved token.json.
        config, err := google.ConfigFromJSON(b, drive.DriveReadonlyScope)
        if err != nil {
                log.Fatalf("Unable to parse client secret file to config: %v", err)
        }
        client := getClient(config)

        srv, err := drive.New(client)
        if err != nil {
                log.Fatalf("Unable to retrieve Drive client: %v", err)
        }

/*
       res, err := srv.Files.Get(
        //"1pjmhD3moj4Rhdm4ilMKRd9kpF68GPO08",
        "1CQ3JdCxFNj-MYOvxYeiGVqYjZNF7gelR",
         
    ).Download()
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    result, _ := ioutil.ReadAll(res.Body)
            _ = ioutil.WriteFile("test.jpeg", result, 0644)
*/
        parentID := "1l-cQFuwP_3mrsVSa-rFGxWDhTv8KH4qa"
        //parentID := "1CQ3JdCxFNj-MYOvxYeiGVqYjZNF7gelR"
        r, err := srv.Files.List().PageSize(1000).
                Fields("nextPageToken, files(id, name)").
                Q("mimeType='image/jpeg' and '"+ parentID + "' in parents").
                //Q("'1CQ3JdCxFNj-MYOvxYeiGVqYjZNF7gelR' in parents and mimeType=image/jpeg").
                //Q("mimeType=image/jpeg").
                Do()
        if err != nil {
                log.Fatalf("Unable to retrieve files: %v", err)
        }
        fmt.Println("Files:")
        if len(r.Files) == 0 {
                fmt.Println("No files found.")
        } else {
                for _, i := range r.Files {
                        fmt.Printf("%s (%s)\n", i.Name, i.Id)
                }
        }

        for _, i := range r.Files {
            res, err := srv.Files.Get(i.Id).Download()
            if err != nil {
                log.Fatalf("Error: %v", err)
            }
            f, err := os.OpenFile("test/"+i.Name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
            if err != nil {
                log.Fatalf("Unable to write image: %v", err)
            }
            result, _ := ioutil.ReadAll(res.Body)
            f.Write(result)
        }
}

