package gormstorage

import (
	"calendar-api/internal/core"
	"context"
	"errors"
	"gorm.io/gorm"
)

func (gs *GormStorage) InstallOrUpdateExtension(ctx context.Context, email string, extensionID uint, additionalData string) error {
	db := gs.db.WithContext(ctx)

	var extensionData core.ExtensionData
	result := db.Model(&core.ExtensionData{}).
		Where(&core.ExtensionData{Extension: extensionID, UserEmail: email}).
		First(&extensionData)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		result = db.Create(&core.ExtensionData{
			UserEmail:      email,
			Extension:      extensionID,
			AdditionalData: additionalData,
		})

		return result.Error
	}

	extensionData.AdditionalData = additionalData
	result = db.Save(&extensionData)

	return result.Error
}

func (gs *GormStorage) DeleteExtension(ctx context.Context, email string, extensionID uint) error {
	db := gs.db.WithContext(ctx)

	result := db.Delete(&core.ExtensionData{Extension: extensionID, UserEmail: email})

	return result.Error
}

func (gs *GormStorage) GetExtensionData(ctx context.Context, email string, extensionID uint) (*core.ExtensionData, error) {
	db := gs.db.WithContext(ctx)

	var extensionData core.ExtensionData
	result := db.Where(&core.ExtensionData{Extension: extensionID, UserEmail: email}).First(&extensionData)

	return &extensionData, result.Error
}
