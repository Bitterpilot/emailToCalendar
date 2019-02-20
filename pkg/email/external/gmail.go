package external

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"

	"github.com/bitterpilot/emailToCalendar/pkg/email"
)

func (srv *gmailSrv) ListEmails(user string) []*email.Msg {
	query := fmt.Sprintf("from:%s subject:%s",
		viper.GetString(`gmailFilter.sender`),
		viper.GetString(`gmailFilter.subject`))
	label := viper.GetString(`gmailFilter.label`)
	listMsg, err := srv.srv.Users.Messages.List(user).LabelIds(label).Q(query).Do()
	if err != nil {
		log.Printf("Could not list emails: %v", err)
	}

	// convert gmail.message struct to email.Msg
	var emails []*email.Msg
	for _, msg := range listMsg.Messages {
		ms := &email.Msg{
			ExternalID: msg.Id,
			ThreadID:   msg.ThreadId,
		}
		emails = append(emails, ms)
	}

	return emails
}

func (srv *gmailSrv) GetEmail(user string, msg *email.Msg) *email.Msg {
	getMsg, err := srv.srv.Users.Messages.Get(user, msg.ExternalID).Do()
	if err != nil {
		log.Printf("Could not retrieve email: %v", err)
	}

	// Check Mime types
	// start by setting a high part number so if a part of the desired mime
	// type is not found we can fail gracefully.
	// parts is an array so a 0 default part num would result in selecting the
	// first element in the array a multi-part email is unlikely to have
	// more than 9000 parts ... right?
	var partNum = 9999
	for i := 0; i < len(getMsg.Payload.Parts); i++ {
		if getMsg.Payload.Parts[i].MimeType == "text/html" {
			partNum = i
		}
	}
	if partNum == 9999 {
		log.Fatalf(
			`
Error: gmail.go/GetMessage text/html Mime Type not found in msgID: %s`,
			msg.ExternalID)
	}

	data := getMsg.Payload.Parts[partNum].Body.Data
	// Confirm Encoding type is base64 before proceding
	partHeaders := getMsg.Payload.Parts[partNum].Headers
	for key := range partHeaders {
		if partHeaders[key].Name == "Content-Transfer-Encoding" {
			switch {
			case partHeaders[key].Value == "base64":
				msg.Body = data
			case partHeaders[key].Value == "quoted-printable":
				msg.Body = data
			default:
				fmt.Printf("can not identify Content-Transfer-Encoding got: %v",
					partHeaders[key].Value)
			}

		}
	}

	// convert gmail.message struct to email.Msg
	msg.ReceivedTime = getMsg.InternalDate

	return msg
}

type gmailSrv struct {
	srv *gmail.Service
}

// NewGmailSrv ...
func NewGmailSrv(user string) email.External {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config, user)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	return &gmailSrv{srv}
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, user string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := fmt.Sprintf("%s/token.json", user)
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		err = os.Mkdir(user, 0700)
		if err != nil {
			log.Fatalf("Unable to create Directory for oauth token: %v", err)
		}
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
