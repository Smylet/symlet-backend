package maintenance

import (
	"database/sql"

	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/common"
)


type HostelMaintenanceRequest struct {
	common.AbstractBaseModel
	HostelID uint `gorm:"not null"`
	Hostel hostel.Hostel

	StudentID uint `gorm:"not null"`
	Student student.Student

	ResolvedByID sql.NullInt16
	ResolvedBy *users.User `gorm:"foreignKey:ResolvedByID"`
	
	RequestImages []HostelMaintenanceRequestImage `gorm:"foreignKey:HostelMaintenanceRequestID"`
	Subject string `gorm:"not null; size:64"`
	Description string `gorm:"not null size:1023"`
	ResolveStatus string `gorm:"default:open; oneof: 'open' 'closed' 'in-progress'"`
	Resolved bool `gorm:"default:false"`


	ResolvedDate sql.NullTime
}

type HostelMaintenanceRequestImage struct {
	common.AbstractBaseImageModel
	HostelMaintenanceRequestID uint
	HostelMaintenanceRequest HostelMaintenanceRequest
}


type WorkOrder struct {
	common.AbstractBaseModel

    HostelMaintenanceRequestID uint
	HostelMaintenanceRequest   HostelMaintenanceRequest

    VendorID                   uint
	Vendor                     users.User
    Description                string `gorm:"size:1023"`

    Status                     string `gorm:"size:16;default:'open'; oneof: 'open' 'pending' 'rejected' 'approved' 'in-progress' 'cancelled' 'completed'"`
	Cost                       float64 `gorm:"default:0"`
	CompletionDate            sql.NullTime
	Comments 				 []WorkOrderComment `gorm:"foreignKey:WorkOrderID"`
	WorkOrderImages			 []WorkOrderImage `gorm:"foreignKey:WorkOrderID"`
}

type WorkOrderImage struct {
	common.AbstractBaseImageModel
	WorkOrderID uint
	WorkOrder WorkOrder
}

type WorkOrderComment struct {
	common.AbstractBaseModel
	WorkOrderID uint
	WorkOrder WorkOrder
	Comment string `gorm:"size:1023"`
	CommentedBy uint
	CommentedByUser users.User `gorm:"foreignKey:CommentedBy"`
}


