package models

import "time"

type RawLog struct {
    Source  string // "ServiceA"
    Message string // "User signed in"
    Time    time.Time
}

type ParsedLog struct {
    Source     string
    Message    string
    Level      string // "INFO", "ERROR", etc.
    Timestamp  time.Time
}


type ErrorLog struct {
 	Source  string // "ServiceA"
    Message string // "User signed in"
    Error error
    Time    time.Time
}

type LogType struct {
	Service string 				 `json:"service"`
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
