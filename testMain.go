package main

import (
	"log"
	"net/http"

	"./cache"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{}
var cacheUser *cache.CacheUser

type Message struct {
	User    string `json:"user"`
	Message string `json:"message"`
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
	r.GET("/ws", handleConnections)
	go handleMessages()
	log.Println("http server started on :8000")
	err := r.Run()
	if err != nil {
		panic(err.Error())
	}
}

func Login(c *gin.Context) {
	person := cache.User{
		Login:    c.PostForm("username"),
		Password: c.PostForm("password"),
	}
	cacheUser.RLock()
	for k, v := range cacheUser.Users {
		if k != person.Login {
			c.JSON(404, "Пользователь не существует")
			cacheUser.RUnlock()
			return
		} else {
			if v != person {
				c.JSON(400, "Неверный логин, или пароль")
				cacheUser.RUnlock()
				return
			} else {
				c.HTML(http.StatusOK, "cabinet.html", nil)
			}
		}
	}
	cacheUser.RUnlock()
}

func handleConnections(c *gin.Context) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
