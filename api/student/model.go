package student

import (
	"time"

	"github.com/Smylet/symlet-backend/api/reference"
	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/common"
	//"github.com/Smylet/symlet-backend/utilities/db"
)

// func init(){
// 	db.RegisterModel(
// 		&Student{})
// }

// Student is a form of user model for our application
type Student struct {
	common.AbstractBaseModel
	User         users.User                    `gorm:"polymorphic:User"`
	UniversityID uint                          `gorm:"not null"`
	University   reference.ReferenceUniversity `gorm:"foreignKey:UniversityID"`
	Department   string                        `gorm:"not null"`
	YearOfEntry  time.Time                       
	ExpectedGraduationYear   time.Time
	StudentIdentificationNumber string		`gorm:"not null"`
	


}
