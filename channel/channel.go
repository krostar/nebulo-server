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
	ID       int       `json:"-" db:"id, primarykey, autoincrement, notnull"`
	Name     string    `json:"name" db:"name, size:42, notnull"`
	Creation time.Time `json:"creation" db:"creation, notnull"`
}
