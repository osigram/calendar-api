package types

type User struct {
	Email          string          `json:"email" gorm:"primaryKey"`
	Name           string          `json:"name"`
	PicturePath    string          `json:"picturePath"`
	ExtensionsData []ExtensionData `json:"-"`
}
