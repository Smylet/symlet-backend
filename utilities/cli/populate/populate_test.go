package populate

import (
	"testing"

	"gorm.io/driver/sqlite"

	"github.com/Smylet/symlet-backend/api/booking"
	"github.com/Smylet/symlet-backend/api/hostel"
	"github.com/Smylet/symlet-backend/api/manager"
	"github.com/Smylet/symlet-backend/api/reference"
	"github.com/Smylet/symlet-backend/api/student"
	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/api/vendor"
	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func TestPopulate(t *testing.T) {
	// Mock the flag values
	mockFlagValue := []string{"amenities", "university"}
	cmd := &cobra.Command{}
	cmd.Flags().StringSliceP("table", "T", mockFlagValue, "Mock flags")

	// Mock the reference models
	mockReferenceModel := &MockReferenceModel{}
	referenceModelMap := map[string]reference.ReferenceModelInterface{
		"amenities":  mockReferenceModel,
		"university": mockReferenceModel,
	}

	// Mock the database
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Execute the Populate function
	err = PopulateReference(cmd, nil, referenceModelMap, db)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify that the Populate method of the mock reference model was called
	if mockReferenceModel.PopulateCalled != 2 {
		t.Errorf("Expected Populate to be called 2 times, but it was called %d times", mockReferenceModel.PopulateCalled)
	}

	// Test the case where a flag doesn't match any reference model
	mockFlagValue = []string{"invalid_table"}
	cmd = &cobra.Command{}
	cmd.Flags().StringSliceP("table", "T", mockFlagValue, "Mock flags")
	err = PopulateReference(cmd, nil, referenceModelMap, db)
	if err == nil || err.Error() != "invalid option invalid_table" {
		t.Errorf("Expected error 'invalid option invalid_table', but got: %v", err)
	}
}

func TestPopulateData(t *testing.T) {
	// Mock the command for populating data
	cmd := &cobra.Command{}

	// Mock the database for populating data
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Execute the DataCommand to populate data
	err = DataCommand.RunE(cmd, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// Assert there is are 3 users in the database, a student, a hostel manager, vendor
	var users []users.User
	db.Find(&users)
	if len(users) != 3 {
		t.Errorf("Expected 3 users in the database, but got %d", len(users))
	}
	var student student.Student
	db.First(&student)
	if student.ID == 0 {
		t.Errorf("Expected a student in the database, but got none")
	}
	var hostelManager manager.HostelManager
	db.First(&hostelManager)
	if hostelManager.ID == 0 {
		t.Errorf("Expected a hostel manager in the database, but got none")
	}
	var vendor vendor.Vendor
	db.First(&vendor)
	if vendor.ID == 0 {
		t.Errorf("Expected a vendor in the database, but got none")
	}

	// Assert an Hostel with a manager of the hostel manager
	var hostel hostel.Hostel
	db.Find("manager_id = ?", hostelManager.ID).First(&hostel)
	if hostel.ID == 0 {
		t.Errorf("Expected a hostel in the database, but got none")
	}

	// Assert a booking with the student and the hostel
	var booking booking.HostelBooking
	db.Find("student_id = ? AND hostel_id = ?", student.ID, hostel.ID).First(&booking)
	if booking.ID == 0 {
		t.Errorf("Expected a booking in the database, but got none")
	}

	// Add additional assertions as needed for the DataCommand branch
}

// MockReferenceModel is a mock implementation of ReferenceModelInterface
type MockReferenceModel struct {
	common.AbstractBaseReferenceModel
	PopulateCalled int
}

func (m *MockReferenceModel) Populate(db *gorm.DB) error {
	m.PopulateCalled++
	return nil
}

func (m *MockReferenceModel) GetTableName() string {
	return "mock_reference_model"
}
