package config

import (
	"os"
)

var (
	Port        string
	Homeserver  string
	Username    string
	Password    string
	Database    string
	RecoveryKey string
	DebugLog    bool
)

func init() {
	// Read env
	Port = os.Getenv("GMA_PORT")
	Homeserver = os.Getenv("GMA_HOMESERVER")
	Username = os.Getenv("GMA_USERNAME")
	Password = os.Getenv("GMA_PASSWORD")
	Database = os.Getenv("GMA_DATABASE")
	RecoveryKey = os.Getenv("GMA_RECOVERYKEY")

	if os.Getenv("GMA_PORT") == "" {
		Port = "8080"
	}

	if os.Getenv("GMA_DATABASE") == "" {
		Database = "/data/gma.db"
	}

	if os.Getenv("GMA_DEBUGLOG") == "true" {
		DebugLog = true
	}

	if Homeserver == "" || Username == "" || Password == "" {
		panic("Please define a Homeserver, Username and Password")
	}
}
