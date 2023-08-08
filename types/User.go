package types

type User struct {
	Email          string          `json:"email"`
	Name           string          `json:"name"`
	PicturePath    string          `json:"picturePath" db:"picture_path"`
	ExtensionsUsed []ExtensionData `json:"-" db:"extension_data"`
}
