# gin_framework_learn

This repo is for learn for gin framework

## GIN 專案設定
- [x] 學習 https://www.youtube.com/watch?v=-IvZBkLh_Lo

### 專案初始化

透過 go mod 初始化一個新專案 --> dependency 會管理在 go.mod

### 安裝 gin 

go get github.com/gin-gonic/gin

### 建立 router 邏輯

設定 main.go 在 package main 如下

```golang=
package main

import "github.com/gin-gonic/gin"

func main () { // entry point
    // create router
    router := gin.Default()
    
    router.Run(":8000") // <--- setup run with port 8000
}
```

### GET/POST route

透過 router 可以設定 GET/POST route

```golang=
...
/** 
 GET /ping
**/
router.GET("/ping", func (c *gin.Context) {
  c.JSON(200, gin.H{
      "message": "ping",
  })
})
/** 
 POST /ping/:id
**/
router.POST("/ping/:id", func (c *gin.Context) {
  id := c.Param("id")
  c.JSON(200, gin.H{
    "id": id,
  })
})
...
```
## 透過 router.Group 實作 User API
- [x] 學習 https://www.youtube.com/watch?v=gaBwPaQjxPY

### 建立 User POJO

```golang==
package pojo

type User struct {
  Id       int `json:"UserId"`
  Name     string `json:"UserName"`
  Password string `json:"UserPassword"`
  Email    string `json:"UserEmail"`
}
```
***Notice*** : json 後面帶入的 key 是 JSON 物件的 key

### 建立 User Service

```golang=
package service

import (
  "web/pojo"
  "net/http"
  "github.com/gin-gonic/gin"
)

var userList = []pojo.User{} // pojo
// GET /users
func FindAllUser(c *gin.Context) {
  c.JSON(http.StatusOK, userList)
}
// POST /users
func PostUser(c *gin.Context) {
  user := pojo.User{}
  // parse context body into pojo
  err := c.BindJSON(&user)
  if err != nil {
    c.JSON(http.StatusNotAcceptable, "Error:" + err.Error())
    return
  }
  userList = append(userList, user)
  c.JSON(http.StatusCreated, "Successfully created")
}
```

### 建立 Group for user route

```golang=
package src

import (
	"web/service"
	"github.com/gin-gonic/gin"
)

func AddUserRouter(r *gin.RouterGroup) {
	user := r.Group("/users")
	user.GET("/", service.FindAllUsers)
	user.POST("/", service.PostUser)
}
```

### 加入 main Router

```golang
package main

import (
	"web/src"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	v1 := router.Group("/v1")
	src.AddUserRouter(v1)
	router.Run(":8000")
}
```

## 建立 PUT/DELETE users

- [x] https://www.youtube.com/watch?v=tHxsLsNRHYs

### 新增 PUT/DELETE users service

```golang=
...
// DELETE User
func DeleteUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	for idx, user := range userList {
		if user.Id == userId {
			userList = append(userList[:idx], userList[idx+1:]...)
			c.JSON(http.StatusOK, "Successfully delete")
			return
		}
	}
	c.JSON(http.StatusNotFound, "Delete Resource not found")
}  
...
// PUT user
func PutUser(c *gin.Context) {
	updatedUser := pojo.User{}
	err := c.BindJSON(&updatedUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, "ParseError")
		return
	}
	userId, _ := strconv.Atoi(c.Param("id"))
	for key, user := range userList {
		if userId == user.Id {
			userList[key] = updatedUser
			log.Println(userList[key])
			c.JSON(http.StatusOK, "Success")
			return
		}
	}
	c.JSON(http.StatusNotFound, "Resource not found")
}
```
### 更新 user route

```golang
package src

import (
	"web/service"

	"github.com/gin-gonic/gin"
)

func AddUserRouter(r *gin.RouterGroup) {
	user := r.Group("/users")
	user.GET("/", service.FindAllUsers)
	user.POST("/", service.PostUser)
	user.DELETE("/:id", service.DeleteUser)
	user.PUT("/:id", service.PutUser)
}
```

## 加入 gorm

- [x] https://www.youtube.com/watch?v=Qe3ekoD_tcw

### 參考 gorm 官網

gorm.io

### 建立 DB 連線

Step 1: 建立 Config

```golang=
package config

type Config struct {
	DBUser     string `json:"DBUser"`
	DBPassword string `json:"DBPassword"`
	DBPort     string `json:"DBPort"`
	DBName     string `json:"DBName"`
	DBHost     string `json:"DBHost"`
	Port       string `json:"Port"`
}
```

Step 2: 加入 dotenv 並且引入 autoload
```shell=
go get github.com/joho/godotenv
```
修改 main 如下
```golang=
package main

import (
	"fmt"
	"log"
	"os"
	"web/config"
	"web/database"
	"web/src"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var Config = config.Config{}

func main() {
	router := gin.Default()
	Config.Port = os.Getenv("PORT")
	Config.DBPort = os.Getenv("DB_PORT")
	Config.DBPassword = os.Getenv("DB_PASSWORD")
	Config.DBUser = os.Getenv("DB_USER")
	Config.DBName = os.Getenv("DB_NAME")
	Config.DBHost = os.Getenv("DB_HOST")
	log.Printf("%v", Config)
	v1 := router.Group("/v1")
	src.AddUserRouter(v1)
	router.Run(fmt.Sprintf(":%s", Config.Port))
}
```
Step 3： 新增 DBConnect.go
```shell=
mkdir database
touch DBConnect.go
```

```golang=
package database

import (
	"fmt"
	"log"
	"web/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBconnect *gorm.DB

var err error

func GetDSN(config *config.Config) string {
	// "user=yuanyu password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Taipei"
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Taipei",
		config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort,
	)
}
func DB(config *config.Config) {
	// https://github.com/go-gorm/postgres
	DBconnect, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  GetDSN(config),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}

```

Step 4: 把 DBConnect 以 goroutine 方式引入

```golang
package main

import (
	"fmt"
	"log"
	"os"
	"web/config"
	"web/database"
	"web/src"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var Config = config.Config{}

func main() {
	router := gin.Default()
	Config.Port = os.Getenv("PORT")
	Config.DBPort = os.Getenv("DB_PORT")
	Config.DBPassword = os.Getenv("DB_PASSWORD")
	Config.DBUser = os.Getenv("DB_USER")
	Config.DBName = os.Getenv("DB_NAME")
	Config.DBHost = os.Getenv("DB_HOST")
	log.Printf("%v", Config)
	go func() {
		database.DB(&Config)
	}()
	v1 := router.Group("/v1")
	src.AddUserRouter(v1)
	router.Run(fmt.Sprintf(":%s", Config.Port))
}

```

Step 5: 把存取邏輯放到 pojo

```golang=
package pojo

import "web/database"

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

```
Step 6: 建立 DB 與 Table

```sql=
CREATE DATABASE Demo;
```
```sql=
CREATE TABLE users (
    id int,
    name varchar(45),
    password varchar(45),
    email varchar(45),
    primary key(id),
);
```
```sql=
INSERT INTO users (id, name, password, email) VALUES 
(1, 'Wilson', '123468', 'Wilson@gmail.com'),
(2, 'tom', '123468', 'tom@gmail.com'),
(3, 'sherry', '123468', 'sherry@gmail.com');
```
STEP 7: 更新 service 與 router

```golang=
...
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
...
```
```golang=
package src

import (
	"web/service"

	"github.com/gin-gonic/gin"
)

func AddUserRouter(r *gin.RouterGroup) {
	user := r.Group("/users")
	user.GET("/", service.FindAllUsers)
	user.GET("/:id", service.FindUserWithId)
	user.POST("/", service.PostUser)
	user.DELETE("/:id", service.DeleteUser)
	user.PUT("/:id", service.PutUser)
}

```
## 新增 POST/DELETE/PUT users

- [x] https://www.youtube.com/watch?v=CO9HqCMmwlY

Step1: 新增 POST/DELETE/PUT users POJO

```golang=
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

```

Step2: 新增 POST/DELETE/PUT users services

```golang=
package service

import (
	"log"
	"net/http"
	"strconv"
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
	c.JSON(http.StatusOK, "Successfuly Delete")
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

```

## 新增 defaultWritter, middleware logger 與 BasicAuth

- [x]  https://www.youtube.com/watch?v=UJfi3ppkqRk

### 設定 gin.DefaultWriter

把 log 寫入一個檔案

更新 main.go 如下

```golang=
...
// setup logger
func setupLogging() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main () {
  setupLogging()
  ...
}
...
```

這樣做把 logger 透過 os.Create 產生一個 fileWriter

先把 os.Stdout 與 fileWriter 包裝成一個 io.Writer

在把 gin.DefaultWriter 指定到這個 io.Writer

這樣一來 gin.Router 接到 logger 除了寫入檔案外, 還有在 os.Stdout 印出執行結果

### 設定 middlewares Logger

透過 middleware 改寫 logger 格式

建立 middlewares/Logger.go

```golang=
package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] %s %s %d \n",
			params.ClientIP,
			params.TimeStamp,
			params.Method,
			params.Path,
			params.StatusCode,
		)
	})
}
```

在 main func 內使用 router.Use 來套用

```golang
...
func main() {
  setupLogging()
  router := gin.Default()
  router.use(middlewares.Logger())
  ...
}
```
### basicAuth

BasicAuth 是指在 request header 做帳密驗證

這邊可以透過 gin.BasicAuth 還有 gin.Account 來設定

使用方式如下

```golang=
func main() {
  ...
    router.Use(gin.BasicAuth(
    gin.Accounts{os.Getenv("BASIC_AUTH_USER"): os.Getenv("BASIC_AUTH_PASSWORD")})
    )
  ...
}
```

## TODO

GORM database migration:

https://gorm.io/docs/migration.html
