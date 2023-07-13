package types

type User struct {
	Email          string      `json:"email"`
	Name           string      `json:"name"`
	PicturePath    string      `json:"photoPath"`
	ExtensionsUsed []Extension `json:"-"`
}
