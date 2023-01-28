package message

import (
	"time"
)

type Message struct {
	Name    string
	Message string
	When    time.Time
}
