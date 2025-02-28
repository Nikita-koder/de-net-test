package model

import (
	"de-net/libs/1_domain_methods/helpers"

	"gorm.io/gorm"
)

func (ut *UserTask) BeforeCreate(tx *gorm.DB) (err error) {
	if ut.ID == "" {
		ut.ID = helpers.GenerateUUID()
	}
	return nil
}
