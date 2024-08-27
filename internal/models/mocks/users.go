package mocks

import (
	"time"

	"github.com/timenglesf/personal-site/internal/models"
)

type UserModel struct{}

func getMockAdminUser() models.User {
	return models.User{
		ID:        "1",
		Username:  "Tim",
		Password:  "password",
		Email:     "tim@email.com",
		IsAdmin:   true,
		CreatedAt: time.Date(2024, 8, 26, 12, 0, 0, 0, time.UTC),
		Posts: []models.Post{
			getMockPost(1), getMockPost(2), getMockPost(3), getMockPost(4),
		},
	}
}

func (m *UserModel) GetAdmin() (*models.User, error) {
	admin := getMockAdminUser()
	return &admin, nil
}

func (m *UserModel) Insert(username, email, password string, isAdmin bool) (string, error) {
	return "2", nil
}

func (m *UserModel) GetByEmail(email string) (*models.User, error) {
	if email == "tim@email.com" {
		admin := getMockAdminUser()
		return &admin, nil
	}
	return nil, models.ErrNoRecord
}

func (m *UserModel) Authenticate(email, password string, checkAdmin bool) (*models.User, error) {
	if email == "tim@email.com" && password == "password" {
		admin := getMockAdminUser()
		return &admin, nil
	}
	return nil, models.ErrInvalidCredentials
}
