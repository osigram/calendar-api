package types

type ExtensionData struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	UserEmail      string `gorm:"references:Email;foreignKey:UserEmail;constraint:OnDelete:CASCADE;"`
	Extension      uint
	AdditionalData string `json:"additionalData"`
}
