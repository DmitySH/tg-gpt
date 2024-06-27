package usecase

import (
	"context"
	"fmt"

	"github.com/DmitySH/tg-gpt/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	beginCommand = "begin"
	endCommand   = "end"
	helpCommand  = "help"
)

type OnCommandUsecaseChatSessionStorage interface {
	HasChatSession(ctx context.Context, userID int64) (bool, error)
	DeleteChatSession(_ context.Context, userID int64) error
	CreateChatSession(ctx context.Context, userID int64, history domain.ChatHistory) error
}

type OnCommandUsecase struct {
	chatSessionStorage OnCommandUsecaseChatSessionStorage
	botAPI             *domain.TGBotAPI
	chatterFabric      func() (domain.Chatter, error)
}

func NewOnCommandUsecase(chatSessionStorage OnCommandUsecaseChatSessionStorage,
	botAPI *domain.TGBotAPI) OnCommandUsecase {
	return OnCommandUsecase{
		chatSessionStorage: chatSessionStorage,
		botAPI:             botAPI,
	}
}

func (u OnCommandUsecase) HandleBotCommand(ctx context.Context, update domain.TGUpdate) error {
	switch update.Message.Command() {
	case beginCommand:
		if err := u.onBeginCommand(ctx, update); err != nil {
			return fmt.Errorf("can't process begin command: %w", err)
		}
	case endCommand:
		if err := u.onEndCommand(ctx, update); err != nil {
			return fmt.Errorf("can't process end command: %w", err)
		}
	case helpCommand:
		if err := u.onHelpCommand(ctx, update); err != nil {
			return fmt.Errorf("can't process help command: %w", err)
		}
	default:
		if err := u.onUnknownCommand(ctx, update); err != nil {
			return fmt.Errorf("can't process unknown command: %w", err)
		}
	}

	return nil
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

	err = u.chatSessionStorage.CreateChatSession(ctx, userID, domain.ChatHistory{})
	if err != nil {
		return fmt.Errorf("can't create new chat session: %w", err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "New chat started")
	_, err = u.botAPI.Send(msg)
	if err != nil {
		return fmt.Errorf("can't send message about chat start: %w", err)
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

	err = u.chatSessionStorage.DeleteChatSession(ctx, userID)
	if err != nil {
		return fmt.Errorf("can't delete chat session: %w", err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Chat ended")
	_, err = u.botAPI.Send(msg)
	if err != nil {
		return fmt.Errorf("can't send message about chat start: %w", err)
	}

	return nil
}

func (u OnCommandUsecase) onHelpCommand(_ context.Context, update domain.TGUpdate) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`/begin - create new chat
/end - end current active chat`,
	)
	_, err := u.botAPI.Send(msg)
	if err != nil {
		return fmt.Errorf("can't send help message: %w", err)
	}

	return nil
}

func (u OnCommandUsecase) onUnknownCommand(_ context.Context, update domain.TGUpdate) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
	_, err := u.botAPI.Send(msg)
	if err != nil {
		return fmt.Errorf("can't send message about unknown command: %w", err)
	}

	return nil
}
