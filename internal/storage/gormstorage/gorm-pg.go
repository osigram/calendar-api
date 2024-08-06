package gormstorage

import (
	types2 "calendar-api/internal/core"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormStorage struct {
	db *gorm.DB
}

func NewStorage(connectionString string) (*GormStorage, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	err = initTable(db, &types2.User{}, err)
	err = initTable(db, &types2.Event{}, err)
	err = initTable(db, &types2.Tag{}, err)
	err = initTable(db, &types2.ExtensionData{}, err)

	return &GormStorage{db}, err
}

func initTable(db *gorm.DB, table any, prevError error) error {
	if prevError == nil {
		return db.AutoMigrate(table)
	}

	return nil
}
