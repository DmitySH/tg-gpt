package gpt

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	client *openai.Client
}

func NewOpenAI(cfg OpenAIConfig) (*OpenAI, error) {
	config := openai.DefaultConfig(cfg.APIKey)
	proxyUrl, err := url.Parse(cfg.ProxyURL)
	if err != nil {
		return nil, err
	}

	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(cfg.ProxyAuth))
	transport := &http.Transport{
		Proxy:              http.ProxyURL(proxyUrl),
		ProxyConnectHeader: http.Header{"Proxy-Authorization": []string{basicAuth}},
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}

	client := openai.NewClientWithConfig(config)

	return &OpenAI{
		client: client,
	}, nil
}

func (o *OpenAI) GetCompletion(ctx context.Context, content string) (string, error) {
	resp, err := o.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("can't create chat completion: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
