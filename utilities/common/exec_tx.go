package common

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

func ExecTx(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	// Begin a new transaction
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	// Execute the function within the transaction context
	if err := fn(tx); err != nil {
		// If the function returned an error, rollback and return the error
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
