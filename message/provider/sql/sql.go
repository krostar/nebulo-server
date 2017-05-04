package sql

import (
	"fmt"
	"math"
	"time"

	"github.com/krostar/nebulo-golib/log"
	gp "github.com/krostar/nebulo-golib/provider"

	"github.com/krostar/nebulo-server/channel"
	"github.com/krostar/nebulo-server/message"
	"github.com/krostar/nebulo-server/message/provider"
	"github.com/krostar/nebulo-server/user"
)

// Provider implements the methods needed to manage messages
// for every SQL based database
type Provider struct {
	*gp.RootProvider
	provider.Provider
}

func (p *Provider) Create(sender user.User, receiver user.User, chann channel.Channel, msg message.SecureMsg) (m *message.Message, err error) {
	m = &message.Message{
		ChannelID:  chann.ID,
		SenderID:   sender.ID,
		ReceiverID: receiver.ID,

		Message:   msg.Message,
		Keys:      msg.Keys,
		Integrity: msg.Integrity,
	}
	if err = p.DB.Create(m).Error; err != nil {
		return nil, fmt.Errorf("unable to insert message: %v", err)
	}

	return m, nil
}

func (p *Provider) List(receiver user.User, chann channel.Channel, lastRead time.Time, limit int) (m []*message.Message, err error) {
	where := &message.Message{
		ChannelID:  chann.ID,
		ReceiverID: receiver.ID,
	}

	if lastRead.IsZero() {
		lastRead = time.Now()
	}
	log.Debugln(lastRead, limit)
	whereTime := "posted"
	if limit < 0 {
		whereTime += " < "
	} else if limit > 0 {
		whereTime += " > "
	} else {
		return nil, nil
	}
	whereTime += "?"
	limit = int(math.Abs(float64(limit)))
	m = make([]*message.Message, limit)
	if err = p.DB.Where(where).Where(whereTime, lastRead).Limit(limit).Find(&m).Error; err != nil {
		return nil, fmt.Errorf("unable to insert message: %v", err)
	}

	for _, mm := range m {
		if err = p.DB.Find(&mm.Channel, mm.ChannelID).Error; err != nil {
			return nil, fmt.Errorf("unable to get channel for message %d: %v", mm.ChannelID, err)
		}
		if err = p.DB.Find(&mm.Sender, mm.SenderID).Error; err != nil {
			return nil, fmt.Errorf("unable to get sender for message %d: %v", mm.SenderID, err)
		}
		if err = p.DB.Find(&mm.Receiver, mm.ReceiverID).Error; err != nil {
			return nil, fmt.Errorf("unable to get receiver for message %d: %v", mm.ReceiverID, err)
		}
	}

	return m, nil
}
