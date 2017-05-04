package channel

import (
	"errors"
	"time"

	"github.com/krostar/nebulo-server/user"
)

var (
	// ErrNotFound is throw when a chan is not found
	ErrNotFound = errors.New("channel not found")
	// ErrNil is throw when a chan is nil
	ErrNil = errors.New("chan is nil")
)

// Channel is the modelisation of a channel
type Channel struct {
	ID        int `json:"-" gorm:"column:id; not null"`
	CreatorID int `json:"-" gorm:"column:creator_id; not null"`

	Name             string      `json:"name" gorm:"column:name; size:64"`
	Created          time.Time   `json:"created" gorm:"column:created; not null" sql:"DEFAULT:current_timestamp"`
	Creator          user.User   `json:"creator" gorm:"ForeignKey:CreatorID; save_associations:false"`
	Members          []user.User `json:"members" gorm:"many2many:channel_memberships"`
	MembersCanEdit   bool        `json:"members_can_edit" gorm:"column:members_can_edit; not null" sql:"DEFAULT:false"`
	MembersCanInvite bool        `json:"members_can_invite" gorm:"column:members_can_invite; not null" sql:"DEFAULT:false"`
}

// UserMembership is the link between a channel and a user
type UserMembership struct {
	ID        int `json:"-" gorm:"column:id; not null"`
	ChannelID int `json:"-" gorm:"column:channel_id; not null"`
	UserID    int `json:"-" gorm:"column:user_id; not null"`

	User    user.User `json:"user" gorm:"ForeignKey:UserID; save_associations:false"`
	Invited time.Time `json:"invited" gorm:"column:invited; not null" sql:"DEFAULT:current_timestamp"`
	Joined  time.Time `json:"joined" gorm:"column:joined"`
}

// TableName is the table name in database
func (um *UserMembership) TableName() string {
	return "channel_memberships"
}
