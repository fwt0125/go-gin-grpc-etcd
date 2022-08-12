package repository

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"user/internal/service"
)

type User struct {
	UserId         uint   `gorm:"primarykey"`
	UserName       string `gorm:"unique"`
	NickName       string
	PasswordDigest string
}

const (
	PasswordCost = 12
)

func (user *User) CheckUserExist(req *service.UserRequest) bool {
	if err := DB.Where("user_name=?", req.UserName).First(&user).Error; err == gorm.ErrRecordNotFound {
		return false
	}
	return true

}

func (user *User) ShowUserInfo(req *service.UserRequest) (err error) {
	if exist := user.CheckUserExist(req); exist {
		return nil
	}
	return errors.New("username not fund")
}

func BuildUser(user User) *service.UserModel {
	userModel := service.UserModel{
		UserId:   uint32(user.UserId),
		UserName: user.UserName,
		NickName: user.NickName,
	}
	return &userModel
}

func (user *User) UserCreate(req *service.UserRequest) error {
	if exist := user.CheckUserExist(req); exist {
		return nil
	}

	userData := User{
		UserName: req.UserName,
		NickName: req.NickName,
	}
	user.SetPassword(req.Password)
	err := DB.Create(&userData).Error
	return err
}

func (user *User) SetPassword(password string) error {
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(fromPassword)
	return nil
}

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}
