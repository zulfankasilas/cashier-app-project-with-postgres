package repository

import (
	"a21hc3NpZ25tZW50/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SessionsRepository struct {
	db *gorm.DB
}

func NewSessionsRepository(db *gorm.DB) SessionsRepository {
	return SessionsRepository{db}
}

func (u *SessionsRepository) AddSessions(session model.Session) error {
	if err := u.db.Create(&session).Error; err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) DeleteSessions(tokenTarget string) error {
	if err := u.db.Where("token = ?", tokenTarget).Delete(&model.Session{}).Error; err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) UpdateSessions(session model.Session) error {
	data := u.db.Model(&model.Session{}).Where("username = ?", session.Username).Updates(&session)

	if data.Error != nil {
		return data.Error
	}

	return nil
}

func (u *SessionsRepository) TokenValidity(token string) (model.Session, error) {
	session, err := u.SessionAvailToken(token)

	if err != nil {
		return model.Session{}, err
	}

	if u.TokenExpired(session) {
		return model.Session{}, fmt.Errorf("Token is Expired!")
	}

	return session, nil
}

func (u *SessionsRepository) SessionAvailName(name string) (model.Session, error) {
	result := model.Session{}

	if err := u.db.Where("username = ?", name).First(&result).Error; err != nil {
		return model.Session{}, err
	}

	return result, nil
}

func (u *SessionsRepository) SessionAvailToken(token string) (model.Session, error) {
	result := model.Session{}

	if err := u.db.Where("token = ?", token).First(&result).Error; err != nil {
		return model.Session{}, err
	}

	return result, nil
}

func (u *SessionsRepository) TokenExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}
