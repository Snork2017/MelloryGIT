package cache

import (
	"errors"
	"sync"
)

type (
	CacheUser struct {
		Users map[string]ValueUser
		sync.RWMutex
	}
	User struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	ValueUser interface {
		GetUser() ValueUser
	}
)

func (user User) GetUser() ValueUser {
	return user
}

func NewUser() *CacheUser {
	users := make(map[string]ValueUser)
	cache := CacheUser{
		Users: users,
	}
	return &cache
}

func (c *CacheUser) SetUser(key string, value ValueUser) {
	c.Lock()
	c.Users[key] = value
	c.Unlock()
}

func (c *CacheUser) DeleteUser(key string) error {
	c.Lock()

	defer c.Unlock()

	if _, found := c.Users[key]; !found {
		return errors.New("Key not found")
	}

	delete(c.Users, key)

	return nil
}
