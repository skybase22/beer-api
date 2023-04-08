package models

import (
	"time"
)

// Log
type Log struct {
	MongoModel `bson:"-"`
	Timestamp  time.Time `json:"timestamp" bson:"timestamp"`
	Message    string    `json:"message" bson:"message"`
}
