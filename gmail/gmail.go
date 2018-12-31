package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "gmail/token.json"
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
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
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

func newService() *gmail.Service {
	b, err := ioutil.ReadFile("gmail/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	return srv
}

func main() {
	srv := newService()
	user := "me"
	msgID := "***REMOVED***"
	r, err := srv.Users.Messages.Get(user, msgID).Format("full").Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Check Mime types
	// start by setting a high part number so if a part of the desiered mime type is not found
	// we can fail gracefully.
	// parts is an array so a 0 default part num would result in selectig the first element in the array
	// a milti-part email is unlikly to have more than 9000 parts (a foolish assumption?)
	var partNum = 9999
	for i := 0; i < len(r.Payload.Parts); i++ {
		if r.Payload.Parts[i].MimeType == "text/html" {
			partNum = i
		}
	}
	if partNum == 9999 {
		log.Fatalf("Mime Type not found")
	}

	// Decode and print
	data := r.Payload.Parts[partNum].Body.Data
	dataType := r.Payload.Parts[partNum].MimeType
	emailBytes, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("%s\n%s", dataType, emailBytes)

}
