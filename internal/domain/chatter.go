package domain

import "context"

type Chatter interface {
	GetCompletion(ctx context.Context, content string) (string, error)
}
