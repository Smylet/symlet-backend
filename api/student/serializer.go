package student

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Smylet/symlet-backend/api/reference"
	"github.com/Smylet/symlet-backend/api/users"
)

type StudentSerializer struct {
	UniversityUID uuid.UUID `json:"university_uid"`

	Department                  string    `json:"department"`
	YearOfEntry                 time.Time `json:"year_of_entry" `
	ExpectedGraduationYear      time.Time `json:"expected_graduation_year"`
	StudentIdentificationNumber string    `json:"student_identification_number"`
	Student                     *Student  `json:"-"`
}

func (h *StudentSerializer) Create(ctx *gin.Context, db *gorm.DB, AWSsession *session.Session) error {
	h.Student = &Student{
		Department:                  h.Department,
		YearOfEntry:                 h.YearOfEntry,
		ExpectedGraduationYear:      h.ExpectedGraduationYear,
		StudentIdentificationNumber: h.StudentIdentificationNumber,
	}

	payload, err := users.GetAuthPayloadFromCtx(ctx)
	if err != nil {
		return fmt.Errorf("unable to retrieve user payload from context %w", err)
	}
	// Does this User already have a Profile?
	err = db.Model(&users.User{}).Preload(clause.Associations).Where("id = ?", payload.UserID).First(&h.Student.User).Error
	if err != nil {
		return fmt.Errorf("unable to retrieve user with id %v %w", payload.UserID, err)
	}
	if h.Student.User.RoleID != 0 {
		return fmt.Errorf("user is already associated with a %v", h.Student.User.RoleType)
	}

	err = db.Model(&reference.ReferenceUniversity{}).Where("uid = ?", h.UniversityUID).First(&h.Student.University).Error
	if err != nil {
		return err
	}
	err = db.Model(&Student{}).Create(&h.Student).Error
	if err != nil {
		return err
	}
	return nil
}

func (h *StudentSerializer) Get(db *gorm.DB, UID uuid.UUID) error {
	err := db.Model(&Student{}).Preload(clause.Associations).Where("uid = ?", UID).First(&h.Student).Error
	if err != nil {
		return err
	}
	return nil
}

// func (h *StudentSerializer) List(db *gorm.DB, queryParams HostelQueryParams) ([]Student, error) {

// }

func (serializer StudentSerializer) Response() interface{} {
	return map[string]interface{}{
		"user": map[string]interface{}{
			"uid":        serializer.Student.User.UID,
			"username":   serializer.Student.User.Username,
			"email":      serializer.Student.User.Email,
			"created_at": serializer.Student.User.CreatedAt,
			"updated_at": serializer.Student.User.UpdatedAt,
		},
		"profile": map[string]interface{}{
			"uid":        serializer.Student.User.Profile.UID,
			"first_name": serializer.Student.User.Profile.FirstName,
			"last_name":  serializer.Student.User.Profile.LastName,
			"bio":        serializer.Student.User.Profile.Bio,
			"image":      serializer.Student.User.Profile.Image,
			"created_at": serializer.Student.User.Profile.CreatedAt,
			"updated_at": serializer.Student.User.Profile.UpdatedAt,
		},
		"department":                    serializer.Student.Department,
		"year_of_entry":                 serializer.Student.YearOfEntry,
		"expected_graduation_year":      serializer.Student.ExpectedGraduationYear,
		"student_identification_number": serializer.Student.StudentIdentificationNumber,
		"created_at":                    serializer.Student.CreatedAt,
		"updated_at":                    serializer.Student.UpdatedAt,
	}

}
