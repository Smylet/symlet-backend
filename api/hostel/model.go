package hostel

import (
	"time"

	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/api/users"

	"github.com/Smylet/symlet-backend/api/core"
	//"github.com/Smylet/symlet-backend/utilities/db"
)

// func init(){
// 	db.RegisterModel(
// 		&Hostel{},
// 		&HostelStudent{},
// 	)
// }

type Hostel struct {
	core.AbstractBaseModel
	Name       string             `gorm:"not null"`
	University string             `gorm:"not null"`
	Address    string             `gorm:"not null"`
	City       string             `gorm:"not null"`
	State      string             `gorm:"not null"`
	Country    string             `gorm:"not null"`
	ManagerID  uint               `gorm:"not null"`
	Manager    users.User         `gorm:"foreignKey:ManagerID"`
	Students   []*student.Student `gorm:"many2many:hostel_students;"`
}

// HostelStudent is the join table between Hostel and Student
// It holds the relationship between the two entities and other metadata
type HostelStudent struct {
	core.AbstractBaseModel
	StudentID     uint
	HostelID      uint
	CheckInDate   time.Time
	CheckOutDate  time.Time
	RoomNumber    string
	CurrentHostel bool

	// Other metadata fields specific to the student-hostel relationship
}
