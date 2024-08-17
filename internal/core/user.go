package core

type User struct {
	Email          string          `json:"email" gorm:"primaryKey"`
	Name           string          `json:"name"`
	Picture        string          `json:"picturePath"`
	ExtensionsData []ExtensionData `json:"-"`
	Sessions       []Session       `json:"sessions"`
}

func (u *User) GetClaims() map[string]interface{} {
	return map[string]interface{}{
		//"id":      float64(u.ID),
		"email":   u.Email,
		"name":    u.Name,
		"picture": u.Picture,
	}
}

func NewUserFromGoogleClaims(claims map[string]interface{}) (user User, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	//idFloat, _ := claims["id"].(float64)
	//id := uint(idFloat)
	email := claims["email"].(string)
	name := claims["name"].(string)
	picture := claims["picture"].(string)

	return User{
		//Model:   gorm.Model{ID: id},
		Email:   email,
		Name:    name,
		Picture: picture,
	}, err
}
