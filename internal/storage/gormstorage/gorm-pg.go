package gormstorage

import (
	"calendar-api/internal/core"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormStorage struct {
	db *gorm.DB
}

func NewStorage(connectionString string) (*GormStorage, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	err = initTable(db, &core.User{}, err)
	err = initTable(db, &core.Event{}, err)
	err = initTable(db, &core.Tag{}, err)
	err = initTable(db, &core.ExtensionData{}, err)

	return &GormStorage{db}, err
}

func initTable(db *gorm.DB, table any, prevError error) error {
	if prevError == nil {
		return db.AutoMigrate(table)
	}

	return nil
}
