package common

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelInterface interface {
	isModel() bool
}

// AbstractBaseModel holds common fields for all tables.

type AbstractBaseModel struct {
	UID       uuid.UUID  `pg:"type:uuid"`
	ID        uint       `json:"id" gorm:"primaryKey" swagger:"description:The primary key for the record."`
	CreatedAt time.Time  `json:"created_at" swagger:"description:Record creation time."`
	UpdatedAt time.Time  `json:"updated_at" swagger:"description:Last record update time."`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index" swagger:"description:Record deletion time. Nil if record is not deleted."`
}

// Before create hook
func (m *AbstractBaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.UID = uuid.New()
	return
}

func (m AbstractBaseModel) isModel() bool {
	return true
}

type AbstractBaseReferenceModel struct {
	AbstractBaseModel
}

type AbstractBaseImageModel struct {
	AbstractBaseModel
	ImageURL string `gorm:"not null"`
}
