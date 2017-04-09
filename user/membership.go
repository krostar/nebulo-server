package user

import "github.com/krostar/nebulo/channel"

// ChannelMembership is the link between a channel and a user
type ChannelMembership struct {
	ID        int `json:"-" gorm:"column:id; not null"`
	ChannelID int `json:"-" gorm:"column:channel_id; not null"`
	UserID    int `json:"-" gorm:"column:user_id; not null"`

	Channel *channel.Channel `json:"channel"`
	User    *User            `json:"user"`
}
