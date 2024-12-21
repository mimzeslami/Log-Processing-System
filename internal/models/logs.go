package models

import "time"

type StructuredLog struct {
	StatusCode int
	API        string
	Message    string
	Timestamp  time.Time
	IPAddress  string
}
