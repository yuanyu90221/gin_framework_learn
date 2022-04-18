package pojo

import (
	"log"
	"web/database"
)

type User struct {
	Id       int    `json:"UserId"`
	Name     string `json:"UserName"`
	Password string `json:"UserPassword"`
	Email    string `json:"UserEmail"`
}

func FindAllUserService() []User {
	var users []User
	database.DBconnect.Find(&users)
	return users
}

func FindByUserId(userId int) User {
	var user User
	database.DBconnect.Where("id = ?", userId).First(&user)
	return user
}

func CreateUser(user User) User {
	database.DBconnect.Create(user)
	return user
}

func DeleteUser(userId int) bool {
	result := database.DBconnect.Where("id = ?", userId).Delete(&User{})
	return result.RowsAffected >= 1
}

func UpdateUser(userId int, user User) bool {
	log.Printf("%v, %v", userId, user)
	result := database.DBconnect.Model(&User{}).Where("id = ?", userId).Updates(user)
	return result.RowsAffected >= 1
}
