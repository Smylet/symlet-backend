package review

import (
	"github.com/Smylet/symlet-backend/api/core"
	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/student"
)


type Review struct{
	core.AbstractBaseModel
	HostelID uint `gorm:"not null"`
	Hostel hostel.Hostel
	Reviews []ReviewType `gorm:"foreignKey:ReviewID references:ID"`
	
}

type ReviewType struct{
	core.AbstractBaseModel
	ReviewName string `gorm:"not null" index:"idx_reviewer_id_review_name oneof: 'security' 'location' 'general' 'amenities' 'water' 'electricity' 'caretaker'"` 
	ReviewerID uint `gorm:"not null" index:"idx_reviewer_id_review_name"`
	Reviewer student.Student `gorm:"foreignKey:ReviewerID"`

	ReviewID uint `gorm:"not null"`
	Review Review

	Rating uint `gorm:"not null; min:1 max:5"` // 1-5 stars
	Description string `gorm:"not null size:1023"`
}