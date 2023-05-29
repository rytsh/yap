package config

var AppName = "yap"

var BuildVars = struct {
	Version string
	Commit  string
	Date    string
}{
	Version: "v0.0.0",
	Commit:  "-",
	Date:    "-",
}
