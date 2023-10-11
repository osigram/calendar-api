package gormstorage

import (
	"calendar-api/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormStorage struct {
	db *gorm.DB
}

func NewStorage(connectionString string) (*GormStorage, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	err = initTable(db, &types.User{}, err)
	err = initTable(db, &types.Event{}, err)
	err = initTable(db, &types.Tag{}, err)
	err = initTable(db, &types.ExtensionData{}, err)

	return &GormStorage{db}, err
}

func initTable(db *gorm.DB, table any, prevError error) error {
	if prevError == nil {
		return db.AutoMigrate(table)
	}

	return nil
}
