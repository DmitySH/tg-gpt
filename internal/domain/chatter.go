package domain

import "context"

type Chatter interface {
	GetCompletion(ctx context.Context, chatHistory ChatHistory, content string) (string, error)
}

type ChatHistory []string
