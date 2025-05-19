package userService

import (
	"crypto/sha256"
	"encoding/hex"
	"message/app/model"
	"message/config/database"
	"time"
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

func GetUserById(id int) (*model.User, error) {
	user := model.User{}
	result := database.DB.Where(
		&model.User{
			Id: id,
		}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return &user, nil
	}
}

func UpdateSignOutTime(id int) error {
	result := database.DB.Model(model.User{}).Where("id = ?", id).Update("sign_out_time", time.Now())
	return result.Error
}
