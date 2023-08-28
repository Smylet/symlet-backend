package review

import (
	"database/sql"

	"github.com/Smylet/symlet-backend/api/core"
	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/api/vendor")


type HostelReview struct{
	core.AbstractBaseModel
	HostelID uint `gorm:"not null"`
	Hostel hostel.Hostel `gorm:"foreignKey:HostelID"`	

    ReviewerID  uint   `gorm:"not null" index:"idx_reviewer_id_review_name"`
    Reviewer    student.Student `gorm:"foreignKey:ReviewerID"`

    Rating      float32 `gorm:"not null check:(rating >= 1) AND (rating <= 5)"`

    // Review types
    SecurityRating   float32 `gorm:"check:(security_rating >= 1) AND (security_rating <= 5)"`
    SecurityComment  sql.NullString  `gorm:"size:127"`
    LocationRating   sql.NullFloat64 `gorm:"check:(location_rating >= 1) AND (location_rating <= 5)"`
    LocationComment sql.NullString  `gorm:"size:127"`
    GeneralRating sql.NullFloat64`gorm:"check:(general_rating >= 1) AND (general_rating <= 5)"`
    GeneralComment   sql.NullString `gorm:"size:127"`
    AmenitiesRating  sql.NullFloat64 `gorm:"check:(amenities_rating >= 1) AND (amenities_rating <= 5)"`
    AmenitiesComment sql.NullString `gorm:"size:127"`
    WaterRating      sql.NullFloat64 `gorm:"check:(water_rating >= 1) AND (water_rating <= 5)"`
    WaterComment     sql.NullString `gorm:"size:127"`
    ElectricityRating sql.NullFloat64 `gorm:"check:(electricity_rating >= 1) AND (electricity_rating <= 5)"`
    ElectricityComment sql.NullString `gorm:"size:127"`
    CaretakerRating  sql.NullFloat64 `gorm:"check:(caretaker_rating >= 1) AND (caretaker_rating <= 5)"`
    CaretakerComment sql.NullString `gorm:"size:127"`
}

type VendorReview struct{
	core.AbstractBaseModel
	VendorID uint `gorm:"not null"`
	Vendor vendor.Vendor `gorm:"foreignKey:VendorID"`

	ReviewerID uint `gorm:"not null"`
	Reviewer users.User `gorm:"foreignKey:ReviewerID"`
	
	Rating          sql.NullFloat64    `gorm:"default:0.0"`//avg of all ratings for this vendor
	GeneralComment   sql.NullString `gorm:"size:127"`
	GeneralRating sql.NullFloat64 
	Quality         sql.NullFloat64    `gorm:"check:(quality >= 1) AND (quality <= 5)"`
	QualityComment  sql.NullString	`gorm:"size:127"`
	Timeliness      sql.NullFloat64    `gorm:"check:(timeliness >= 1) AND (timeliness <= 5)"`
	TimelinessComment sql.NullString `gorm:"size:127"`
	Communication   sql.NullFloat64    `gorm:"check:(communication >= 1) AND (communication <= 5)"`
	CommunicationComment sql.NullString `gorm:"size:127"`
	Professionalism sql.NullFloat64    `gorm:"check:(professionalism >= 1) AND (professionalism <= 5)"`
	ProfessionalismComment sql.NullString `gorm:"size:127"`
	CostEffectiveness sql.NullFloat64    `gorm:"check:(cost_effectiveness >= 1) AND (cost_effectiveness <= 5)"`
	CostEffectivenessComment sql.NullString	`gorm:"size:127"`
	Reliability     sql.NullFloat64    `gorm:"check:(reliability >= 1) AND (reliability <= 5)"`
	ReliabilityComment sql.NullString `gorm:"size:127"`
	ProblemSolving  sql.NullFloat64    `gorm:"check:(problem_solving >= 1) AND (problem_solving <= 5)"`
	ProblemSolvingComment sql.NullString `gorm:"size:127"`
	Flexibility     uint    `gorm:"check:(flexibility >= 1) AND (flexibility <= 5)"`
	FlexibilityComment sql.NullString `gorm:"size:127"`
	CustomerSatisfaction sql.NullFloat64    `gorm:"check:(customer_satisfaction >= 1) AND (customer_satisfaction <= 5)"`
	CustomerSatisfactionComment sql.NullString `gorm:"size:127"`
	ResponseTime    sql.NullFloat64    `gorm:"check:(response_time >= 1) AND (response_time <= 5)"`
	ResponseTimeComment sql.NullString `gorm:"size:127"`
	ProblemResolution sql.NullFloat64    `gorm:"check:(problem_resolution >= 1) AND (problem_resolution <= 5)"`
	ProblemResolutionComment sql.NullString `gorm:"size:127"`
	}
	

type HostelManagerReview struct {
    core.AbstractBaseModel
    ReviewerID  uint               `gorm:"not null"`
    Reviewer    users.User         `gorm:"foreignKey:ReviewerID"`

    ManagerID   uint               `gorm:"not null"`
    Manager     users.User         `gorm:"foreignKey:ManagerID"`

    Rating      float32            `gorm:"default:0.0;check:(rating >= 0) AND (rating <= 5)"`
    Description string             `gorm:"not null;size:1023"`

    // Review criteria fields and comments
    CommunicationRating sql.NullFloat64     `gorm:"check:(communication_rating >= 1) AND (communication_rating <= 5)"`
    CommunicationComment sql.NullString    `gorm:"size:127"`
    ProfessionalismRating sql.NullFloat64   `gorm:"check:(professionalism_rating >= 1) AND (professionalism_rating <= 5)"`
    ProfessionalismComment sql.NullString  `gorm:"size:127"`
    ResponsivenessRating sql.NullFloat64   `gorm:"check:(responsiveness_rating >= 1) AND (responsiveness_rating <= 5)"`
    ResponsivenessComment sql.NullString  `gorm:"size:127"`
}
