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

func (ss *Session) CreateForUser(user *models.User) (string, error) {
	token := uuid.NewV4().String()

	session := &models.Session{
		Token:          token,
		ExpirationTime: time.Now().Add(time.Hour),
		UserId:         &user.ID,
	}

	err := ss.DB.Create(session).Error
	return token, err
}

func (ss *Session) CreateForClient(node *models.Client) (string, error) {
	token := uuid.NewV4().String()

	session := &models.Session{
		Token:          token,
		ExpirationTime: time.Now().Add(time.Hour),
		ClientId:       &node.ID,
	}

	err := ss.DB.Create(session).Error
	return token, err
}

func (ss *Session) ValidateUser(token string) (models.User, error) {
	session := &models.Session{Token: token}
	err := ss.DB.Preload("User").Where(session).Find(session).Error

	//validate this session is for a user
	if err == nil && session.UserId == nil {
		err = gorm.ErrRecordNotFound
	}
	return session.User, err
}

func (ss *Session) ValidateClient(token string) (models.Client, error) {
	session := &models.Session{Token: token}
	err := ss.DB.Preload("Client").Preload("User").Preload("Server").Where(session).Find(session).Error

	//validate this session is for a client
	if err == nil && session.ClientId == nil {
		err = gorm.ErrRecordNotFound
	}
	return session.Client, err
}

func (ss *Session) ValidateNode(token string) (models.Node, error) {
	node := &models.Node{Secret: token}
	err := ss.DB.Where(node).First(node).Error
	return *node, err
}

func (ss *Session) Expire(token string) error {
	session := &models.Session{Token: token}
	return ss.DB.Delete(session).Error
}
