package notification


import (
	"github.com/Smylet/symlet-backend/api/core"
)


type Notification struct {
	core.AbstractBaseModel
    UserID      uint   
    Content     string `gorm:"not null size:1023"`
    IsRead      bool   `gorm:"default:false"`
	
	// Type of action associated with the notification (e.g., "maintenance_request", "assignment", etc.)
    ActionType  string 
    ActionID    uint   // ID of the action related to the notification

}

