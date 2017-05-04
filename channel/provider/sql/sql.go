package sql

import (
	"fmt"
	"strings"

	gp "github.com/krostar/nebulo-golib/provider"

	"github.com/krostar/nebulo-server/channel"
	"github.com/krostar/nebulo-server/channel/provider"
	"github.com/krostar/nebulo-server/user"
)

// Provider implements the methods needed to manage channels
// for every SQL based database
type Provider struct {
	*gp.RootProvider
	provider.Provider
}

// Create create a channel if needed, or return an exsting one with the same requirements
func (p *Provider) Create(name string, creator user.User, members []user.User) (c *channel.Channel, err error) {
	toFind := channel.Channel{
		Name:      name,
		CreatorID: creator.ID,
	}

	c, err = p.Find(toFind)
	if err != nil && err == channel.ErrNotFound {
		c = &channel.Channel{
			Name:      name,
			CreatorID: creator.ID,
			Members:   members,
		}
		if err = p.DB.Create(c).Error; err != nil {
			return nil, fmt.Errorf("unable to insert channel: %v", err)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("unable to find channel: %v", err)
	}

	return c, nil
}

// List return a list of channel
func (p *Provider) List(u user.User, offset int, limit int) (list map[string]*channel.Channel, err error) {
	list = make(map[string]*channel.Channel)
	// TODO: its ugly as hell, done it in a hurry, need to redo, thats bad, some cats dies because of me, see with gorm team?
	// TODO: please whoever read this, forgive me, I was tired, I'M NOT HERE TO SUFFER OKKAAY?

	type fakeChannel struct {
		channel.Channel
		Members string
	}
	tmp := []fakeChannel{}
	if err = p.DB.Table("`channels` AS `c`").Select("`c`.*, GROUP_CONCAT(`um`.id) AS `members`").
		Joins("INNER JOIN `channel_memberships` AS `cc` ON `cc`.`channel_id` = `c`.`id`").
		Joins("INNER JOIN `channel_memberships` AS `cm` ON `cm`.`channel_id` = `c`.`id`").
		Joins("INNER JOIN `users` AS `um` ON `cm`.`user_id` = `um`.`id`").
		Where("`cc`.`user_id` = ?", u.ID).Group("`c`.`id`").Limit(limit).Offset(offset).Scan(&tmp).Error; err != nil {
		return nil, fmt.Errorf("unable to get channels list for user %d: %v", u.ID, err)
	}

	for _, c := range tmp {
		var chnel channel.Channel
		chnel = c.Channel
		chnel.Members = []user.User{}
		membersID := strings.Split(c.Members, ",")
		if err = p.DB.Where(membersID).Find(&chnel.Members).Error; err != nil {
			return nil, fmt.Errorf("unable to get channels list for user %d: %v", u.ID, err)
		}
		if err = p.DB.Where(chnel.CreatorID).Find(&chnel.Creator).Error; err != nil {
			return nil, fmt.Errorf("unable to get channels list for user %d: %v", u.ID, err)
		}
		list[chnel.Name] = &chnel
	}

	return list, nil
}
