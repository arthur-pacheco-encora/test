package debug

import (
	"encoding/json"
	"errors"
	"fmt"
	stdlog "log"
	"time"
)

type log struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

var (
	database = []log{}
	debug    = false
)

func Init(isActive bool) {
	debug = isActive
	database = nil
}

func NewMessage(message string) {
	stdlog.Println(message)
	if debug {
		now := time.Now() //nolint:forbidigo
		database = append(database, log{Message: message, Time: now.Format(time.RFC3339)})
	}
}

func GetAllMessages() []log {
	return database
}

func GetJson() ([]byte, error) {
	json, err := json.Marshal(GetAllMessages())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to encode the debug data: %v", err))
	}
	return json, nil
}
