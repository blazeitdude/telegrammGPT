package historyCache

import (
	"sync"
)

type UserCache struct {
	Messages []Message
	mu       sync.RWMutex
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewUserCache() *UserCache {
	return &UserCache{
		Messages: make([]Message, 0),
	}
}

func (c *UserCache) AddMessage(message Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Messages = append(c.Messages, message)
}

func (c *UserCache) GetMessages() []Message {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]Message(nil), c.Messages...)
}

type Cache struct {
	users map[string]*UserCache
	mu    sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		users: make(map[string]*UserCache),
	}
}

func (c *Cache) GetUserCache(userID string) *UserCache {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, found := c.users[userID]; !found {
		c.users[userID] = NewUserCache()
	}
	return c.users[userID]
}
