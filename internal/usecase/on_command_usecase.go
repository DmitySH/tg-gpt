package usecase

import (
	"context"
	"fmt"

	"github.com/DmitySH/tg-gpt/internal/domain"
	"github.com/DmitySH/tg-gpt/internal/pkg/loggy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	beginCommand = "begin"
	endCommand   = "end"
)

type ChatSessionStorage interface {
	HasChatSession(ctx context.Context, userID int64) (bool, error)
	GetChatSession(ctx context.Context, userID int64) (domain.Chatter, error)
	CreateChatSession(ctx context.Context, userID int64, chatter domain.Chatter) error
}

type OnCommandUsecase struct {
	chatSessionStorage ChatSessionStorage
	botAPI             *domain.TGBotAPI
}

func NewOnCommandUsecase(chatSessionStorage ChatSessionStorage, botAPI *domain.TGBotAPI) OnCommandUsecase {
	return OnCommandUsecase{
		chatSessionStorage: chatSessionStorage,
		botAPI:             botAPI,
	}
}

func (u OnCommandUsecase) HandleBotCommand(ctx context.Context, update domain.TGUpdate) {
	var err error
	switch update.Message.Command() {
	case beginCommand:
		err = u.onBeginCommand(ctx, update)
		if err != nil {
			err = fmt.Errorf("can't process begin command: %w", err)
		}
	case endCommand:
		err = u.onEndCommand(ctx, update)
		if err != nil {
			err = fmt.Errorf("can't process end command: %w", err)
		}
	default:
		err = u.onUnknownCommand(ctx, update)
		if err != nil {
			err = fmt.Errorf("can't process unknown command: %w", err)
		}
	}

	if err != nil {
		loggy.Errorln(err)
	}
}

func (u OnCommandUsecase) onBeginCommand(ctx context.Context, update domain.TGUpdate) error {
	userID := update.Message.From.ID
	hasChat, err := u.chatSessionStorage.HasChatSession(ctx, userID)
	if err != nil {
		return fmt.Errorf("can't check if chat exists: %w", err)
	}

	if hasChat {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You already have active chat. To begin new, end previous chat using /end command")
		_, err = u.botAPI.Send(msg)
		if err != nil {
			return fmt.Errorf("can't send message about active chat: %w", err)
		}

		return nil
	}

	return nil
}

func (u OnCommandUsecase) onEndCommand(ctx context.Context, update domain.TGUpdate) error {
	userID := update.Message.From.ID

	hasChat, err := u.chatSessionStorage.HasChatSession(ctx, userID)
	if err != nil {
		return fmt.Errorf("can't check if chat exists: %w", err)
	}

	if !hasChat {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You don't have any active chats. To begin new, use /begin command")
		_, err = u.botAPI.Send(msg)
		if err != nil {
			return fmt.Errorf("can't send message about active chat: %w", err)
		}

		return nil
	}

	return nil
}

func (u OnCommandUsecase) onUnknownCommand(_ context.Context, update domain.TGUpdate) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
	_, err := u.botAPI.Send(msg)
	if err != nil {
		return fmt.Errorf("can't send message about: %w", err)
	}

	return nil
}
