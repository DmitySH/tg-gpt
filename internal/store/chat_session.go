package store

import (
	"context"
	"time"

	"github.com/DmitySH/tg-gpt/internal/domain"
	"github.com/jellydator/ttlcache/v3"
)

type ChatSessionStorage struct {
	store *ttlcache.Cache[int64, domain.ChatHistory]
}

func NewChatSessionStorage(defaultTTL time.Duration, maxCapacity uint64) *ChatSessionStorage {
	return &ChatSessionStorage{
		store: ttlcache.New[int64, domain.ChatHistory](
			ttlcache.WithTTL[int64, domain.ChatHistory](defaultTTL),
			ttlcache.WithCapacity[int64, domain.ChatHistory](maxCapacity),
		),
	}
}

func (c *ChatSessionStorage) GetChatSession(_ context.Context, userID int64) (domain.ChatHistory, error) {
	return c.store.Get(userID).Value(), nil
}

func (c *ChatSessionStorage) HasChatSession(_ context.Context, userID int64) (bool, error) {
	return c.store.Has(userID), nil
}

func (c *ChatSessionStorage) CreateChatSession(_ context.Context, userID int64, history domain.ChatHistory) error {
	c.store.Set(userID, history, ttlcache.DefaultTTL)
	return nil
}

func (c *ChatSessionStorage) DeleteChatSession(_ context.Context, userID int64) error {
	c.store.Delete(userID)
	return nil
}
