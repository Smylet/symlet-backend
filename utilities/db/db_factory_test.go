package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMakeDBProvider(t *testing.T) {
	tests := []struct {
		name              string
		dsn               string
		expectedDialector string
	}{
		{
			name:              "WithSqliteURI",
			dsn:               "sqlite:///tmp/smylet.db",
			expectedDialector: "sqlite",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DB = nil
			db, err := MakeDBProvider(
				tt.dsn,
				time.Second*2,
				2,
				false, // ensure SQLite does not attempt a reset
			)
			assert.Nil(t, err)

			if db != nil { // Only proceed if db is not nil
				assert.NotNil(t, db)
				assert.Equal(t, tt.expectedDialector, db.GormDB().Dialector.Name())

				// expecting the global 'DB' not to be set
				assert.Nil(t, DB)
				assert.Nil(t, db.Close())
			}
		})
	}
}
