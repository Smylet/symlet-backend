package core

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelInterface interface {
	isModel() bool
}



type AbstractBaseModel struct {
	gorm.Model
	UID uuid.UUID `pg:"type:uuid"`
}

//Before create hook
func (m *AbstractBaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.UID = uuid.New()
	return
}

func (m *AbstractBaseModel) isModel() bool {
	return true
}

type AbstractBaseReferenceModel struct{
	AbstractBaseModel
	Name        string `gorm:"size:1023;uniqueIndex"`
    Value       int16  `gorm:"unique"`
    Slug        string `gorm:"size:140;uniqueIndex"`
    Description string `gorm:"size:1023;index"`
}

func (r *AbstractBaseReferenceModel) isReferenceModel() bool {
	return true
}


