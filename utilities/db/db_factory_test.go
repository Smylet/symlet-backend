package db

// func TestMakeDBProvider(t *testing.T) {
// 	config, err := utils.LoadConfig()
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("failed to load config")
// 	}

// 	tests := []struct {
// 		name              string
// 		dsn               string
// 		expectedDialector string
// 	}{
// 		{
// 			name:              "WithSqliteURI",
// 			dsn:               config.DatabaseURI,
// 			expectedDialector: "sqlite",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			db, err := MakeDBProvider(config)
// 			assert.Nil(t, err)

// 			if db != nil { // Only proceed if db is not nil
// 				assert.NotNil(t, db)
// 				assert.Equal(t, tt.expectedDialector, db.GormDB().Dialector.Name())

// 				// expecting the global 'DB' not to be set
// 				assert.Nil(t, db.GormDB())
// 				assert.Nil(t, db.Close())
// 			}
// 		})
// 	}
// }
