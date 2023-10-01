package fixtures

import (
	"github.com/Smylet/symlet-backend/api/users"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// baseFixtures represents base fixtures object.
type baseFixtures struct {
	db *gorm.DB
}

// UnloadFixtures cleans database from the old data.
func (f baseFixtures) UnloadFixtures() error {
	for _, table := range []interface{}{
		users.User{},
		users.Profile{},
		users.VerificationEmail{},
		users.Session{},
	} {
		if err := f.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(table).Error; err != nil {
			return errors.Wrap(err, "error deleting data")
		}
	}
	return nil
}
