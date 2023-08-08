package postgres

import (
	"calendar-api/types"
	"errors"
)

func (s *Storage) GetUser(email string) (*types.User, error) {
	var user types.User
	err := s.db.Get(
		&user,
		`SELECT * from users where email = $1`,
		email,
	)
	if err != nil {
		return nil, errors.New("user not found: " + err.Error())
	}
	err = s.db.Select(
		&(user.ExtensionsUsed),
		`select id, additional_data from extension_data where email = $1`,
		email,
	)
	if err != nil {
		return nil, errors.New("user.extension_data getting from db error")
	}

	return &user, nil
}

func (s *Storage) AddUser(user *types.User) error {
	_, err := s.db.NamedExec(
		`INSERT INTO users (email, name, picture_path) 
		VALUES (:email, :name, :picture_path)`,
		user,
	)

	return err
}
