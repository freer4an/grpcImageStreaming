package repository

import (
	"log/slog"
	"sync"
)

type ImageCache struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewImageCache() *ImageCache {
	return &ImageCache{
		data: make(map[string][]byte),
	}
}

func (c *ImageCache) SaveImage(uid string, image []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[uid] = image
}

func (c *ImageCache) GetImage(uid string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RLock()
	if v, ok := c.data[uid]; ok {
		return v, true
	}
	return nil, false
}

func (c *ImageCache) DeleteImage(uid string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.data[uid]; ok {
		slog.Info("Image deleted",
			slog.String("uid", uid))
	}
	delete(c.data, uid)
}

func (c *ImageCache) ListImages() map[string][]byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data
}
