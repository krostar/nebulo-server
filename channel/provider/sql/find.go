package sql

import (
	"fmt"
	"strings"

	"github.com/krostar/nebulo-server/channel"
	"github.com/krostar/nebulo-server/user"
)

// Find is used to find a channel from the setted field
func (p *Provider) Find(toFind channel.Channel) (c *channel.Channel, err error) {
	c = new(channel.Channel)

	if p.DB.Where(toFind).First(c).RecordNotFound() {
		return nil, channel.ErrNotFound
	}
	if err = p.DB.Error; err != nil {
		return nil, fmt.Errorf("unable to select channel in db: %v", err)
	}

	return c, nil
}

// FindByName is used to find a channel from his name
func (p *Provider) FindByName(u user.User, name string) (c *channel.Channel, err error) {
	c = new(channel.Channel)

	type fakeChannel struct {
		channel.Channel
		Members string
	}
	tmp := fakeChannel{}
	if err = p.DB.Table("`channels` AS `c`").Select("`c`.*, GROUP_CONCAT(`um`.id) AS `members`").
		Joins("INNER JOIN `channel_memberships` AS `cc` ON `cc`.`channel_id` = `c`.`id`").
		Joins("INNER JOIN `channel_memberships` AS `cm` ON `cm`.`channel_id` = `c`.`id`").
		Joins("INNER JOIN `users` AS `um` ON `cm`.`user_id` = `um`.`id`").
		Where("`cc`.`user_id` = ?", u.ID).Where("`c`.`name` = ?", name).Group("`c`.`id`").Limit(1).Scan(&tmp).Error; err != nil {
		return nil, fmt.Errorf("unable to get channels list for user %d: %v", u.ID, err)
	}

	c = &tmp.Channel
	c.Members = []user.User{}
	membersID := strings.Split(tmp.Members, ",")
	if err = p.DB.Where(membersID).Find(&c.Members).Error; err != nil {
		return nil, fmt.Errorf("unable to get channels list for user %d: %v", u.ID, err)
	}
	if err = p.DB.Where(c.CreatorID).Find(&c.Creator).Error; err != nil {
		return nil, fmt.Errorf("unable to get channels list for user %d: %v", u.ID, err)
	}
	return c, nil
}
