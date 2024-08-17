package core

type Session struct {
	ID           uint   `json:"id"`
	UserEmail    string `json:"userEmail"`
	RefreshToken string `json:"-"`
	DeviceData   string `json:"deviceData"`
}

func (session *Session) GetClaims() map[string]interface{} {
	return map[string]interface{}{
		"id":           session.ID,
		"userEmail":    session.UserEmail,
		"deviceData":   session.DeviceData,
		"refreshToken": session.RefreshToken,
	}
}

func NewSessionFromClaims(claims map[string]interface{}) (session Session, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	idFloat, _ := claims["id"].(float64)
	id := uint(idFloat)
	userEmail := claims["userEmail"].(string)
	refreshToken := claims["refreshToken"].(string)
	deviceData := claims["deviceData"].(string)

	return Session{
		ID:           id,
		UserEmail:    userEmail,
		RefreshToken: refreshToken,
		DeviceData:   deviceData,
	}, err
}
