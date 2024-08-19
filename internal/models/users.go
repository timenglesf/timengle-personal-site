package models

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        string     `gorm:"primaryKey"`
	Username  string     `gorm:"unique;not null"`
	Password  string     `gorm:"not null"`
	Email     string     `gorm:"unique;not null:index"`
	IsAdmin   bool       `gorm:"not null;default:false"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	LastLogin *time.Time `gorm:""`
	Posts     []Post     `gorm:"foreignKey:AuthorID"`
}

type UserModel struct {
	DB *gorm.DB
}

func (m *UserModel) GetAdmin() (*User, error) {
	var user User
	// Get single user with is_admin = true
	err := m.DB.Where("is_admin = ?", true).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) Insert(username, email, password string, isAdmin bool) (string, error) {
	admin, err := m.GetAdmin()
	if err != nil {
		if !errors.Is(err, ErrNoRecord) {
			return "", err
		}
	}
	// If admin already exists and we are trying to create another admin
	if isAdmin && admin != nil {
		return "", ErrDuplicateAdmin
	}
	// Create a new user
	user := User{
		ID:        uuid.New().String(),
		Username:  username,
		Email:     strings.TrimSpace(strings.ToLower(email)),
		Password:  password,
		IsAdmin:   isAdmin,
		CreatedAt: time.Now().UTC(),
	}
	err = m.DB.Create(&user).Error
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	var user User
	err := m.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) Authenticate(email, password string, checkAdmin bool) (*User, error) {
	var user User
	cleanEmail := strings.TrimSpace(strings.ToLower(email))
	err := m.DB.Where("email = ?", cleanEmail).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if checkAdmin && !user.IsAdmin {
		return nil, ErrInvalidCredentials
	}
	// validate password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Update LastLogin time
	now := time.Now().UTC()
	user.LastLogin = &now
	if err = m.DB.Save(&user).Error; err != nil {
		return nil, ErrUpdateUser
	}
	return &user, nil
}
