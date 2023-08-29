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
	print("Before create hook")
	m.UID = uuid.New()
	return
}

func (m AbstractBaseModel) isModel() bool {
	return true
}

type AbstractBaseReferenceModel struct{
	AbstractBaseModel
}


type AbstractBaseImageModel struct {
	AbstractBaseModel
	ImageURL string `gorm:"not null"`
}

