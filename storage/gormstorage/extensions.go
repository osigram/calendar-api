package gormstorage

import (
	"calendar-api/types"
	"errors"
	"gorm.io/gorm"
)

func (gs *GormStorage) InstallOrUpdateExtension(email string, extensionID uint, additionalData string) error {
	db := gs.db

	var extensionData types.ExtensionData
	result := db.Model(&types.ExtensionData{}).
		Where(&types.ExtensionData{Extension: extensionID, UserEmail: email}).
		First(&extensionData)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		result = db.Create(&types.ExtensionData{
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

	result := db.Delete(&types.ExtensionData{Extension: extensionID, UserEmail: email})

	return result.Error
}
