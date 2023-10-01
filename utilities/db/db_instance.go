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
