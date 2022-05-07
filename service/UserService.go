package service

import (
	"log"
	"net/http"
	"strconv"
	"web/middlewares"
	"web/pojo"

	"github.com/gin-gonic/gin"
)

// Get User
func FindAllUsers(c *gin.Context) {
	// c.JSON(http.StatusOK, userList)
	users := pojo.FindAllUserService()
	c.JSON(http.StatusOK, users)
}

func FindUserWithId(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user := pojo.FindByUserId(id)
	if user.Id == 0 {
		c.JSON(http.StatusNotFound, "Not found")
		return
	}
	log.Println("User ->", user)
	c.JSON(http.StatusOK, user)
}

// Post User
func PostUser(c *gin.Context) {
	user := pojo.User{}
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, "Error:"+err.Error())
		return
	}
	// userList = append(userList, user)
	newUser := pojo.CreateUser(user)
	c.JSON(http.StatusCreated, newUser)
}

// Delete User
func DeleteUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	isDeleted := pojo.DeleteUser(userId)
	if !isDeleted {
		c.JSON(http.StatusNotFound, "Delete Resource not found")
		return
	}
	c.JSON(http.StatusOK, "Successfully Delete")
}

func PutUser(c *gin.Context) {
	updatedUser := pojo.User{}
	err := c.BindJSON(&updatedUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, "ParseError")
		return
	}
	userId, _ := strconv.Atoi(c.Param("id"))
	isUpdated := pojo.UpdateUser(userId, updatedUser)
	log.Printf("%v", isUpdated)
	if !isUpdated {
		c.JSON(http.StatusNotFound, "Update resource not found")
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

// Login User
func LoginUser(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	user := pojo.CheckUserPassword(name, password)
	if user.Id == 0 {
		c.JSON(http.StatusNotFound, "Error")
		return
	}
	middlewares.SaveSession(c, user.Id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login Successfully",
		"User":    user,
		"Session": middlewares.GetSession(c),
	})
}

// Logout Users
func LogoutUser(c *gin.Context) {
	middlewares.ClearSession(c)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout Successfully",
	})
}

// CheckUserSession
func CheckUserSession(c *gin.Context) {
	sessionId := middlewares.GetSession(c)
	if sessionId == 0 {
		c.JSON(http.StatusUnauthorized, "Error")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Check Session Successfully",
		"User":    middlewares.GetSession(c),
	})
}
