package main

import (
	"fmt"
	"net/http"

	"./cache"
	"github.com/gin-gonic/gin"
)

var cacheUser *cache.CacheUser

func Login(c *gin.Context) {
	person := cache.User{
		Login:    c.PostForm("username"),
		Password: c.PostForm("password"),
	}
	for k, v := range cacheUser.Users {
		if k == person.Login {
			if v != person {
				c.JSON(400, "Неверный логин, или пароль")
				return
			} else {
				c.JSON(200, fmt.Sprintf("Добро пожаловать %v", person.Login))
			}
		} else {
			c.JSON(404, "Пользователь не существует")
			return
		}
	}
}

func main() {
	cacheUser = cache.NewUser()
	newUser := cache.User{
		Login:    "Mellory",
		Password: "123",
	}
	cacheUser.SetUser("Mellory", newUser)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "testMain.html", nil)
	})
	r.POST("/login", Login)
	err := r.Run()
	if err != nil {
		fmt.Println("Ошибка при запуске сервера", err.Error())
	}
}
