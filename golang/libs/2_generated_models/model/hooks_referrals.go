package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ut *Referral) BeforeCreate(tx *gorm.DB) (err error) {
	if ut.ID == "" {
		ut.ID = uuid.New().String()
	}
	return nil
}
