package db

import (
	"io"

	"gorm.io/gorm"
)

// DBProvider is the interface to access the DB.
type DBProvider interface {
	GormDB() *gorm.DB
	Dsn() string
	Close() error
	Reset() error
}

// DBInstance is the base concrete type for DbProvider.
type DBInstance struct {
	*gorm.DB
	dsn     string
	closers []io.Closer
}

// Close will invoke the closers.
func (db *DBInstance) Close() error {
	for _, c := range db.closers {
		err := c.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// Dsn will return the dsn string.
func (db *DBInstance) Dsn() string {
	return db.dsn
}

// Db will return the gorm DB.
func (db *DBInstance) GormDB() *gorm.DB {
	return db.DB
}

// createDefaultExperiment will create the default experiment if needed.
// func createDefaultExperiment(defaultArtifactRoot string, db DBProvider) error {
// 	if tx := db.GormDB().First(&Experiment{}, 0); tx.Error != nil {
// 		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
// 			log.Info("Creating default experiment")
// 			var id int32 = 0
// 			ts := time.Now().UTC().UnixMilli()
// 			exp := Experiment{
// 				ID:             &id,
// 				Name:           "Default",
// 				LifecycleStage: LifecycleStageActive,
// 				CreationTime: sql.NullInt64{
// 					Int64: ts,
// 					Valid: true,
// 				},
// 				LastUpdateTime: sql.NullInt64{
// 					Int64: ts,
// 					Valid: true,
// 				},
// 			}
// 			if tx := db.GormDB().Create(&exp); tx.Error != nil {
// 				return fmt.Errorf("error creating default experiment: %s", tx.Error)
// 			}

// 			exp.ArtifactLocation = fmt.Sprintf("%s/%d", strings.TrimRight(defaultArtifactRoot, "/"), *exp.ID)
// 			if tx := db.GormDB().Model(&exp).Update("ArtifactLocation", exp.ArtifactLocation); tx.Error != nil {
// 				return fmt.Errorf("error updating artifact_location for experiment '%s': %s", exp.Name, tx.Error)
// 			}
// 		} else {
// 			return fmt.Errorf("unable to find default experiment: %s", tx.Error)
// 		}
// 	}
// 	return nil
// }
