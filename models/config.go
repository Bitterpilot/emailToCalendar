package models

// Config holds settings
type Config struct {
	User      string
	Label     string
	Sender    string
	Subject   string
	Calendar  string
	Locations map[string]string
}
