package historyCache

import (
	"sync"
	"telegrammGPT/pkg/gptClient"
)

type UserCache struct {
	Messages []gptClient.Message
	mu       sync.RWMutex
}

func NewUserCache() *UserCache {
	return &UserCache{
		Messages: make([]gptClient.Message, 0),
	}
}

func (c *UserCache) AddMessage(message gptClient.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Messages = append(c.Messages, message)
}

func (c *UserCache) GetMessages() []gptClient.Message {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]gptClient.Message(nil), c.Messages...)
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
