package core

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelInterface interface{
	isModel() bool
}


type AbstractBaseModel struct{
	gorm.Model
	UID uuid.UUID `json:"uid"`
}

func (m AbstractBaseModel)isModel() bool{
	return true
}