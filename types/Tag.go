package types

type Tag struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	EventID uint   `gorm:"constraint:OnDelete:CASCADE;"`
	TagText string `json:"tagText"`
}
