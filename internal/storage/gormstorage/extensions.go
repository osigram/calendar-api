package gormstorage

import (
	"calendar-api/internal/core"
	"errors"
	"gorm.io/gorm"
)

func (gs *GormStorage) InstallOrUpdateExtension(email string, extensionID uint, additionalData string) error {
	db := gs.db

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

func (gs *GormStorage) DeleteExtension(email string, extensionID uint) error {
	db := gs.db

	result := db.Delete(&core.ExtensionData{Extension: extensionID, UserEmail: email})

	return result.Error
}

func (gs *GormStorage) GetExtensionData(email string, extensionID uint) (*core.ExtensionData, error) {
	db := gs.db

	var extensionData core.ExtensionData
	result := db.Where(&core.ExtensionData{Extension: extensionID, UserEmail: email}).First(&extensionData)

	return &extensionData, result.Error
}
