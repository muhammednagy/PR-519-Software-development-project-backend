package chat

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/models"
)


//Connect connect user to user channels on redis
func Connect(rdb *redis.Client, name string) (*models.User, error) {
	if _, err := rdb.SAdd(models.UsersKey, name).Result(); err != nil {
		return nil, err
	}

	u := &models.User{
		Nickname:             name,
		StopListenerChan: make(chan struct{}),
		MessageChan:      make(chan redis.Message),
	}

	if err := u.Connect(rdb); err != nil {
		return nil, err
	}

	return u, nil
}

func Chat(rdb *redis.Client, channel string, content string) error {
	return rdb.Publish(channel, content).Err()
}

func List(rdb *redis.Client) ([]string, error) {
	return rdb.SMembers(models.UsersKey).Result()
}

func GetChannels(rdb *redis.Client, username string) ([]string, error) {

	if !rdb.SIsMember(models.UsersKey, username).Val() {
		return nil, errors.New("user not exists")
	}

	var c []string

	c1, err := rdb.SMembers(models.ChannelsKey).Result()
	if err != nil {
		return nil, err
	}
	c = append(c, c1...)

	// get all user channels (from DB) and start subscribe
	c2, err := rdb.SMembers(fmt.Sprintf(models.UserChannelFmt, username)).Result()
	if err != nil {
		return nil, err
	}
	c = append(c, c2...)
	return c, nil
}