package models

import (
	"errors"
	"fmt"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

const (
	// used to track users that used chat. mainly for listing users in the /users api, in real world chat app
	// such user list should be separated into user management module.
	UsersKey       = "users"
	UserChannelFmt = "user:%s:channels"
	ChannelsKey    = "channels"
)

type User struct {
	ID               uint32             `gorm:"primary_key;auto_increment" json:"id"`
	Nickname         string             `gorm:"size:255;not null;unique" json:"nickname"`
	Email            string             `gorm:"size:100;not null;unique" json:"email"`
	Password         string             `gorm:"size:100;not null;" json:"password"`
	CreatedAt        time.Time          `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time          `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	ChannelsHandler  *redis.PubSub      `gorm:"-"`
	StopListenerChan chan struct{}      `gorm:"-"`
	Listening        bool               `gorm:"-"`
	MessageChan      chan redis.Message `gorm:"-"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Nickname == "" {
			return errors.New("required Nickname")
		}
		if u.Password == "" {
			return errors.New("required Password")
		}
		if u.Email == "" {
			return errors.New("required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid Email")
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("required Password")
		}
		if u.Email == "" {
			return errors.New("required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid Email")
		}
		return nil

	default:
		if u.Nickname == "" {
			return errors.New("required Nickname")
		}
		if u.Password == "" {
			return errors.New("required Password")
		}
		if u.Email == "" {
			return errors.New("required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid Email")
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {

	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"nickname":  u.Nickname,
			"email":     u.Email,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (u *User) Subscribe(rdb *redis.Client, channel string) error {

	userChannelsKey := fmt.Sprintf(UserChannelFmt, u.Nickname)

	if rdb.SIsMember(userChannelsKey, channel).Val() {
		return nil
	}
	if err := rdb.SAdd(userChannelsKey, channel).Err(); err != nil {
		return err
	}

	return u.Connect(rdb)
}

func (u *User) Unsubscribe(rdb *redis.Client, channel string) error {

	userChannelsKey := fmt.Sprintf(UserChannelFmt, u.Nickname)

	if !rdb.SIsMember(userChannelsKey, channel).Val() {
		return nil
	}
	if err := rdb.SRem(userChannelsKey, channel).Err(); err != nil {
		return err
	}

	return u.Connect(rdb)
}

func (u *User) Connect(rdb *redis.Client) error {

	var c []string

	c1, err := rdb.SMembers(ChannelsKey).Result()
	if err != nil {
		return err
	}
	c = append(c, c1...)

	// get all user channels (from DB) and start subscribe
	c2, err := rdb.SMembers(fmt.Sprintf(UserChannelFmt, u.Nickname)).Result()
	if err != nil {
		return err
	}
	c = append(c, c2...)

	if len(c) == 0 {
		fmt.Println("no channels to connect to for user: ", u.Nickname)
		return nil
	}

	if u.ChannelsHandler != nil {
		if err := u.ChannelsHandler.Unsubscribe(); err != nil {
			return err
		}
		if err := u.ChannelsHandler.Close(); err != nil {
			return err
		}
	}
	if u.Listening {
		u.StopListenerChan <- struct{}{}
	}

	return u.doConnect(rdb, c...)
}

func (u *User) doConnect(rdb *redis.Client, channels ...string) error {
	// subscribe all channels in one request
	pubSub := rdb.Subscribe(channels...)
	// keep channel handler to be used in unsubscribe
	u.ChannelsHandler = pubSub

	// The Listener
	go func() {
		u.Listening = true
		fmt.Println("starting the listener for user:", u.Nickname, "on channels:", channels)
		for {
			select {
			case msg, ok := <-pubSub.Channel():
				if !ok {
					return
				}
				u.MessageChan <- *msg

			case <-u.StopListenerChan:
				fmt.Println("stopping the listener for user:", u.Nickname)
				return
			}
		}
	}()
	return nil
}

func (u *User) Disconnect() error {
	if u.ChannelsHandler != nil {
		if err := u.ChannelsHandler.Unsubscribe(); err != nil {
			return err
		}
		if err := u.ChannelsHandler.Close(); err != nil {
			return err
		}
	}
	if u.Listening {
		u.StopListenerChan <- struct{}{}
	}

	close(u.MessageChan)

	return nil
}
