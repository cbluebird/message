package userService

import (
	"crypto/sha256"
	"encoding/hex"
	"message/app/model"
	"message/config/database"
)

func Login(email, password string) bool {
	h := sha256.New()
	h.Write([]byte(password))
	pass := hex.EncodeToString(h.Sum(nil))
	user := model.User{}
	result := database.DB.Where(
		model.User{
			Email:    email,
			Password: pass,
		}).First(&user)
	return result.Error == nil
}

func GetUserByEmail(email string) (*model.User, error) {
	user := model.User{}
	result := database.DB.Where(
		&model.User{
			Email: email,
		}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return &user, nil
	}
}

func CreateUser(user *model.User) error {
	h := sha256.New()
	h.Write([]byte(user.Password))
	pass := hex.EncodeToString(h.Sum(nil))
	user.Password = pass
	return database.DB.Create(user).Error
}
