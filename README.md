# Gmail to gCal
Take a email of shifts from riteq and put them into google calandar

### Sample
>Your schedule for 24 Dec 2018 through to 6 Jan 2019 is shown below
>
>| Day | Date   | Start Work | End Work | Total Hours | Breaks | Pay   | Org Level                           |
>|-----|--------|------------|----------|-------------|--------|-------|-------------------------------------|
>| Mon | 24 Dec | 12:45      | 18:15    | 05:30       | 00:30  | 05:00 | AAAA\Dry Operations\Snr CSO         |
>| Thu | 27 Dec | 13:45      | 20:00    | 06:15       | 00:30  | 05:45 | AAAA\Dry Operations\Snr CSO         |
>| Fri | 28 Dec | 13:45      | 20:00    | 06:15       | 00:30  | 05:45 | AAAA\Dry Operations\Snr CSO         |
>| Sat | 29 Dec | 06:00      | 12:15    | 06:15       | 00:00  | 06:15 | AAAA\Dry Operations\Dry Ops Officer |
>| Sun | 30 Dec | 13:00      | 18:15    | 05:15       | 00:00  | 05:15 | AAAA\Dry Operations\Snr CSO         |
>| Thu | 03 Jan | 13:15      | 21:15    | 08:00       | 00:00  | 08:00 | AAAA\Dry Operations\Dry Ops Officer |
>| Fri | 04 Jan | 07:30      | 14:00    | 06:30       | 00:30  | 06:00 | AAAA\Dry Operations\Snr CSO         |

## Road map
### v0.5
Restructure code to a clean code architecture
```go
.  
├── cmd
│   ├── cli  // full cli version of the app(expects user interaction) calls app.Run()
│   └── faas // functions as a service version(expects events (i.e. mail received, timed event) to trigger
│            // function)
├── pkg  
│   ├── email // gets emails
│   │   │     // expects: User && (MsgID || ThdID || Watch Response)
│   │   │     // returns: err && (email body || MsgID || ThdID || time received)
│   │   │     // methods:
│   │   │           // NewService calls a provider specific function to open a new authentication session.
│   │   │           // Although tempting to have the authentication tokens combine permissions for email 
│   │   │           // and calendar, the email NewService is distinct from the calendar one to avoid the 
│   │   │           // two packages from being tightly coupled.
│   │   │           // It also has the benefit of avoiding complexity if the user wants to map between providers.
│   │   │           type ServiceHandler struct {
│   │   │               user string
│   │   │               providerService interface
│   │   │               logger *log.logger
│   │   │               db *sql.db
│   │   │           }
│   │   │           NewService(user, provider) (svc, err) {
│   │   │               switch {
│   │   │                  case provider: provider methods 
│   │   │                  default: errors.New("must have provider")
│   │   │               }
│   │   │           }
│   │   │           
│   │   │           type provider interface {
│   │   │               (*NewService) Push(bool) (pushRsp, err){if true{watch} if false{stop}}
│   │   │               (*NewService) HandlePush(pushRsp) (historyID, err)
│   │   │               (*NewService) ListRecentEmails(historyID)
│   │   │               (*NewService) ListAllEmails()
│   │   │               (*NewService) GetEmail(MsgID) *email
│   │   │               (*NewService) GetEmails(ThdID || []MsgID) *[]email
│   │   │           }
│   │   ├──  email.go  // implements newService() and provider interface
│   │   └──  providers
│   │        ├──  google.go     // implements the functions listed in provider interface for Gmail
│   │        └──  microsoft.go // implements the functions listed in provider interface for Microsoft
│   │
│   ├── calendar // gets, removes and publishes events
│   │   │        // expects: User && (Event || eventIds)
│   │   │        // returns: Err && (eventIDs || Events)
│   │   │        // methods:
│   │   │           // NewService calls a provider specific function to open a new authentication session.
│   │   │           // Although tempting to have the authentication tokens combine permissions for calendar 
│   │   │           // and email, the calendar NewService is distinct from the email one to avoid the 
│   │   │           // two packages from being tightly coupled.
│   │   │           // It also has the benefit of avoiding complexity if the user wants to map between providers.
│   │   │           type ServiceHandler struct {
│   │   │               user string
│   │   │               providerService interface
│   │   │               logger *log.logger
│   │   │           }
│   │   │           NewService(user, provider) (svc, err) 
│   │   │                
│   │   │           (*NewService) GetEvents(startDate, endDate) ([]event, err)
│   │   │           // maybe these bulk actions can take []event or ...event and if event < 1 call the multiple 
│   │   │           // operation method.
│   │   │           (*NewService) RemoveEvent(eventID) err
│   │   │           (*NewService) RemoveEvents([]eventIDs) err
│   │   │           (*NewService) PublishEvent(event) err
│   │   │           (*NewService) PublishEvents([]event) err
│   │   │           (*NewService) ChangeEvent([]event) err
│   │   ├──  calendar.go  // implements newService() and provider interface
│   │   └──  providers
│   │        ├──  google.go    // implements the functions listed in provider interface for Google
│   │        └──  microsoft.go // implements the functions listed in provider interface for Microsoft
│   │
│   └── app // builds events from email body and checks for existing events
│       │       // expects: User && (shift || email)
│       │       // returns: Err && (event)
│       │       // methods: 
│       │           type email struct {
│       │               IntId    int    // Internal Storage ID
│       │               ExtID    string // ID from the service provider
│       │               ExtThdID string // ID linking emails together from the service provider
│       │               Body     string // base64 encoded
│       │           }
│       │           type shift struct {
│       │               IntId      int    // Internal Storage ID
│       │               ExtEventID string // ID from the service provider
│       │               ExtMsgID   string // ID linking emails together from the service provider
│       │               Comment    string // base64 encoded
│       │               URL        string // direct link to source email
│       │           
│       │           }
│       │           type handler struct {
│       │               logger *log.logger
│       │               db     *sql.db
│       │           }
│       │           NewHandler(logger, db) handler
│       │           RunCli() // function handles all the code to run the cli version expects parameters handed
│       │                    // at runtime
│       │           (*handler) ProcessEmail(MsgID) (*email)
│       │           (*handler) dateRange(body)
│       │           (*handler) checkEncoding(body)
│       │           (*handler) processTable(body) [][]string
│       │           (*handler) tagReader
│       │           (*handler) convertTableLineToShift([]string) shift
│       │           (*handler) convertShiftToEvent(shift) event
│       │           (*handler) compareEvent(event, event) (bool, diff)
│       │           (*handler) convertDates(string) time
│       │           (*handler) 
│       ├── store
│       │   └── db
│       ├── app.go   // designed in a way that it will be easy to break up into individual files such as email.go
│       └── run.go   // contains method that are called from main.go
│
├── vendor
│   └── ... 
├── README.md
├── .gitignore
├── dockerfile
├── go.mod
├── continues.integration
└── ...

```
### v0.6
Recognize a email that changes previously published shifts.
This will probably use the dates in the first line of the email so it will also
need to recognizing the difference between republish and an email for the same
dates but different location.
### v0.7
Functions to batch requests for gmail and calendar