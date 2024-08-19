package models

import (
	"gorm.io/gorm"
)

type Meta struct {
	gorm.Model
	Version     string `gorm:"not null"`
	Name        string `gorm:"not null"`
	LastUpdated string `gorm:"not null"`
	Description string
	Author      string
	Environment string
	BuildNumber string
	License     string
}

type MetaModel struct {
	DB *gorm.DB
}

func (m *MetaModel) InsertMeta(md Meta) error {
	result := m.DB.Create(&md)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *MetaModel) GetMostRecentMeta() (*Meta, error) {
	var meta Meta
	result := m.DB.Order("id desc").First(&meta)
	if result.Error != nil {
		return nil, result.Error
	}
	return &meta, nil
}
