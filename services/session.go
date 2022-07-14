package services

import (
	"github.com/pufferpanel/pufferpanel/v3/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type Session struct {
	DB *gorm.DB
}

func (ss *Session) Create(user *models.User) (string, error) {
	token := uuid.NewV4().String()

	session := &models.Session{
		Token:          token,
		ExpirationTime: time.Now().Add(time.Hour),
		UserId:         user.ID,
	}

	err := ss.DB.Create(session).Error
	return token, err
}

func (ss *Session) Validate(token string) (uint, error) {
	session := &models.Session{Token: token}
	err := ss.DB.Preload("User").Where(session).Find(session).Error
	return session.UserId, err
}

func (ss *Session) Expire(token string) error {
	session := &models.Session{Token: token}
	return ss.DB.Delete(session).Error
}
