package message

import (
	"time"

	"github.com/krostar/nebulo-server/channel"
	"github.com/krostar/nebulo-server/user"
)

type SecureMsg struct {
	Message   []byte `json:"message"`
	Keys      []byte `json:"keys"`
	Integrity []byte `json:"integrity"`
}

type Message struct {
	ID         int `json:"-" gorm:"column:id; not null"`
	ChannelID  int `json:"-" gorm:"column:channel_id; not null"`
	SenderID   int `json:"-" gorm:"column:sender_id; not null"`
	ReceiverID int `json:"-" gorm:"column:receiver_id; not null"`

	Message   []byte `json:"message" gorm:"column:message; type:blob; not null"`
	Keys      []byte `json:"keys" gorm:"column:keys; size:256; not null"`
	Integrity []byte `json:"integrity" gorm:"column:integrity; size:32; not null"`

	Channel  channel.Channel `json:"channel" gorm:"column:ForeignKey:ChannelID; save_associations:false"`
	Sender   user.User       `json:"sender" gorm:"column:ForeignKey:SenderID; save_associations:false"`
	Receiver user.User       `json:"receiver" gorm:"column:ForeignKey:ReceiverID; save_associations:false"`
	Posted   time.Time       `json:"posted" gorm:"column:posted; not null" sql:"DEFAULT:current_timestamp"`
	Seen     time.Time       `json:"seen" gorm:"column:seen" sql:"DEFAULT:NULL"`
}
