package auth

import (
	e "errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/guemidiborhane/factorydigitale.tech/internal/errors"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/movies/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetupModels(database *gorm.DB) {
	db = database
	if !fiber.IsChild() {
		go db.AutoMigrate(&User{})
	}
}

type User struct {
	ID       int    `json:"id" binding:"required"`
	Username string `json:"username" gorm:"uniqueIndex" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`

	LockedAt  time.Time `json:"locked_at"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Authenticable struct {
	Users []User `gorm:"polymorphic:Authenticable"`
}

type UserTracks struct {
	LoginCount uint        `json:"login_count"`
	IPs        [][2]string `json:"ips"`
}

func (user *User) Create() error {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Unexpected(err.Error())
	}

	user.Password = string(password)
	if err := db.Create(&user).Error; err != nil {
		return errors.Unexpected(err.Error())
	}

	return nil
}

func (user *User) Get() error {
	if err := db.First(&user, "id = ? OR username = ?", user.ID, user.Username).Error; err != nil {
		if e.Is(err, gorm.ErrRecordNotFound) {
			return errors.EntityNotFound(err.Error())
		}

		return errors.Unexpected(err.Error())
	}

	return nil
}

func (u *User) Favourites() ([]int, error) {
	var ids []int
	if err := db.Model(models.Favourite{}).Where("user_id = ?", u.ID).Pluck("movie_id", &ids).Error; err != nil {
		return []int{}, err
	}

	return ids, nil
}

type UserResponse struct {
	ID       int    `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

func (user *User) AsJSON() *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
}
