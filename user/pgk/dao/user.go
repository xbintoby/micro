package dao

import (
	"context"
	"log"
	"time"
)

type User struct {
	Uid       int64     `gorm:"column:uid; PRIMARY_KEY"`
	Username  string    `gorm:"column:username"`
	Fullname  string    `gorm:"column:fullname"`
	Password  string    `gorm:"column:password"`
	Birthdata int64     `gorm:"column:birthdata"`
	Bio       string    `gorm:"column:bio"`
	Token     string    `gorm:"column:token"`
	Create_at time.Time `gorm:"column:create_at"`
	Update_at time.Time `gorm:"column:update_at"`
}

func (u User) TableName() string {

	return "user"
}

func Save(ctx context.Context, user *User) {
	err := DB().Create(user)
	if err != nil {
		log.Println("insert fail : ", err)
	}
}
func Update(ctx context.Context, uid int64, u User) {

	user := User{}
	DB().Where("uid = ?", uid).Take(&user)
	//DB().Model(&user).Update("username", "update100")
	DB().Model(&user).Updates(User{
		Username: u.Username,
		Fullname: u.Fullname,
		Token:    u.Token,
	})
}

func Delete(ctx context.Context, uid int64) {
	user := User{}
	DB().Where("uid = ?", uid).Take(&user)
	err := DB().Delete(&user).Error
	if err != nil {
		log.Println("delete fail :", err)
	}
}
func GetInfo(ctx context.Context, uid int64) *User {
	user := User{}
	err := DB().Where("uid = ?", uid).Take(&user).Error

	if err != nil {
		log.Println("select one fail :", err)
	}
	return &user
}
func Login(ctx context.Context, username, password string) *User {
	user := User{}
	err := DB().Where("username = ? and password = ?", username, password).Take(&user).Error

	if err != nil {
		log.Println("select one fail :", err)
	}
	return &user
}
func GetList(ctx context.Context, page int, pagesize int) []User {

	var users []User
	offset := (page - 1) * pagesize
	err := DB().Offset(offset).Limit(pagesize).Find(&users).Error

	if err != nil {
		log.Println("select one fail :", err)
	}
	return users
}

func Count(ctx context.Context) int64 {
	var total int64 = 0
	DB().Model(User{}).Count(&total)
	return total
}
