package types

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Tag struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	EventID uint   `gorm:"constraint:OnDelete:CASCADE;"`
	TagText string `json:"tagText"`
}

func (tag *Tag) Validate() error {
	return validation.ValidateStruct(tag,
		validation.Field(&tag.TagText, validation.Required),
	)
}
