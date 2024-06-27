package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/DmitySH/tg-gpt/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type OnBasicMessageUsecaseChatSessionStorage interface {
	GetChatSession(_ context.Context, userID int64) (domain.ChatHistory, error)
	HasChatSession(ctx context.Context, userID int64) (bool, error)
	CreateChatSession(ctx context.Context, userID int64, history domain.ChatHistory) error
}

type OnBasicMessageUsecase struct {
	chatSessionStorage OnBasicMessageUsecaseChatSessionStorage
	botAPI             *domain.TGBotAPI
	chatter            domain.Chatter
}

func NewOnBasicMessageUsecase(chatSessionStorage OnBasicMessageUsecaseChatSessionStorage,
	botAPI *domain.TGBotAPI, chatter domain.Chatter) OnBasicMessageUsecase {
	return OnBasicMessageUsecase{
		chatSessionStorage: chatSessionStorage,
		botAPI:             botAPI,
		chatter:            chatter,
	}
}

func (u OnBasicMessageUsecase) HandleBasicMessage(ctx context.Context, update domain.TGUpdate) error {
	userID := update.Message.From.ID

	hasChat, err := u.chatSessionStorage.HasChatSession(ctx, userID)
	if err != nil {
		return fmt.Errorf("can't check if chat exists: %w", err)
	}

	if !hasChat {
		return nil
	}

	history, err := u.chatSessionStorage.GetChatSession(ctx, userID)
	if err != nil {
		return fmt.Errorf("can't check if chat exists: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	completion, err := u.chatter.GetCompletion(ctx, history, update.Message.Text)
	if err != nil {
		return fmt.Errorf("can't get completion from chatter: %w", err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, completion)
	_, err = u.botAPI.Send(msg)
	if err != nil {
		return fmt.Errorf("can't send completion message: %w", err)
	}

	history = append(history, update.Message.Text)
	err = u.chatSessionStorage.CreateChatSession(ctx, userID, history)
	if err != nil {
		return fmt.Errorf("can't save new chat session: %w", err)
	}

	return nil
}
