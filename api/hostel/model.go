package hostel

import (
	"time"

	"github.com/The-CuriousX/project/api/student"
	users "github.com/The-CuriousX/project/api/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Hostel struct {
	gorm.Model
	UID uuid.UUID `gorm:"not null"`
	Name string `gorm:"not null"`
	University string `gorm:"not null"`
	Address string `gorm:"not null"`
	City string `gorm:"not null"`
	State string `gorm:"not null"`
	Country string `gorm:"not null"`
	ManagerID uint `gorm:"not null"`
	Manager users.User `gorm:"foreignKey:ManagerID"`
	Students []*student.Student `gorm:"many2many:hostel_students;"`
	
}


// HostelStudent is the join table between Hostel and Student
// It holds the relationship between the two entities and other metadata
type HostelStudent struct {
	gorm.Model
	StudentID     uint
	HostelID      uint
	CheckInDate   time.Time
	CheckOutDate  time.Time
	RoomNumber    string
	CurrentHostel bool

	// Other metadata fields specific to the student-hostel relationship
}

