package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ut *UserTask) BeforeCreate(tx *gorm.DB) (err error) {
	if ut.ID == "" {
		ut.ID = uuid.New().String()
	}
	return nil
}
