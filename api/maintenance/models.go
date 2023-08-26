package maintenance

import (
	"time"

	"github.com/Smylet/symlet-backend/api/core"
	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/api/users"
)


type HostelMaintenanceRequest struct {
	core.AbstractBaseModel
	HostelID uint `gorm:"not null"`
	Hostel hostel.Hostel

	StudentID uint `gorm:"not null"`
	Student student.Student

	ResolvedBy uint
	ResolvedByUser users.User
	
	RequestImages []string `gorm:"not null"`
	Subject string `gorm:"not null; size:64"`
	Description string `gorm:"not null size:1023"`
	ResolveStatus string `gorm:"default:open; oneof: 'open' 'closed' 'in-progress'"`
	Resolved bool `gorm:"not null"`


	ResolvedDate time.Time
}


type WorkOrder struct {
	core.AbstractBaseModel

    HostelMaintenanceRequestID uint
	HostelMaintenanceRequest   HostelMaintenanceRequest

    VendorID                   uint
	Vendor                     users.User
    Description                string `gorm:"size:1023"`

    Status                     string `gorm:"size:16;default:'open' oneof: 'open' 'pending' 'rejected' 'approved' 'in-progress' 'cancelled' 'completed'"`
	Image 					 string `gorm:"size:1023"`
	Cost                       float64
	CompletionDate             time.Time
	Comments 				 []WorkOrderComment `gorm:"foreignKey:WorkOrderID"`
}

type WorkOrderComment struct {
	core.AbstractBaseModel
	WorkOrderID uint
	WorkOrder WorkOrder
	Comment string `gorm:"size:1023"`
	CommentedBy uint
	CommentedByUser users.User `gorm:"foreignKey:CommentedBy"`
}


