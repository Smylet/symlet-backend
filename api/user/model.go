// models.go
package users

import "time"

type User struct {
	ID           uint `gorm:"primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `sql:"not null"`
	Profile      Profile
}

type Profile struct {
	ID       uint `gorm:"primary_key"`
	UserID   uint
	Username string
	Bio      string
	Image    string
}

// func SaveUserWithProfile(user User, profile Profile) error {
// 	tx := db.Begin()

// 	if err := tx.Create(&user).Error; err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	if err := tx.Create(&profile).Error; err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	return tx.Commit().Error
// }
