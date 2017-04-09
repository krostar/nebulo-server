package channel

import (
	"errors"
	"time"
)

var (
	// ErrNotFound is throw when a chan is not found
	ErrNotFound = errors.New("channel not found")
	// ErrNil is throw when a chan is nil
	ErrNil = errors.New("chan is nil")
)

// Channel is the modelisation of a channel
type Channel struct {
	ID int `json:"-" gorm:"column:id; not null"`

	Name    string    `json:"name" gorm:"column:name; size:42; not null"`
	Created time.Time `json:"created" gorm:"column:created; not null"`
}
