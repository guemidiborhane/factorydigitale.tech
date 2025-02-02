package auth

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"

	"github.com/guemidiborhane/factorydigitale.tech/internal/errors"
	"github.com/guemidiborhane/factorydigitale.tech/internal/storage"
	"github.com/guemidiborhane/factorydigitale.tech/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func GetCurrentUser(c *fiber.Ctx) (User, error) {
	user := User{
		ID: c.Locals(USER_ID).(int),
	}

	if err := user.Get(); err != nil {
		return User{}, errors.Unauthorized
	}

	return user, nil
}

func trackingStorageKey(id int) string {
	return fmt.Sprintf("users:%d", id)
}

func (user *User) GetTracks() (UserTracks, error) {
	var storage_key string = trackingStorageKey(user.ID)
	var data UserTracks

	currentData, _ := storage.Redis.Get(storage_key)
	if currentData != nil {
		if err := json.Unmarshal(currentData, &data); err != nil {
			return UserTracks{}, errors.Unexpected(err.Error())
		}
	}

	return data, nil
}

func (user *User) WriteTracks(data *UserTracks) error {
	var storage_key string = trackingStorageKey(user.ID)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.Unexpected(err.Error())
	}

	if err := storage.Redis.Set(storage_key, jsonData, 0); err != nil {
		return errors.Unexpected(err.Error())
	}

	return nil
}

func (user *User) Track(c *fiber.Ctx) error {
	var (
		ip     string = c.IP()
		lastIp string = ""
	)

	storageData, err := user.GetTracks()
	if err != nil {
		return err
	}

	if len(storageData.IPs) > 0 {
		lastIp = storageData.IPs[len(storageData.IPs)-1][0]
	}

	if lastIp != ip {
		storageData.IPs = append(storageData.IPs, [2]string{ip, time.Now().Format(time.RFC3339)})
		user.WriteTracks(&storageData)
	}

	return nil
}

func VerifyUserSession(c *fiber.Ctx) error {
	session, err := storage.Session.Get(c)
	if err != nil {
		return errors.Unauthorized
	}

	if session.Get(AUTH_KEY) == nil {
		return errors.Unauthorized
	}

	userID := session.Get(USER_ID).(int)

	user := &User{
		ID: userID,
	}

	if err := user.Get(); err != nil {
		return errors.Unauthorized
	}

	c.Locals(USER_ID, userID)

	if err := user.Track(c); err != nil {
		utils.WriteToStderr(err)
	}

	return nil
}
