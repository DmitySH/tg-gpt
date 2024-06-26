package usecase

import (
	"context"

	"github.com/DmitySH/tg-gpt/internal/domain"
)

type OnBasicMessageUsecase struct {
	chatter domain.Chatter
}

func NewOnBasicMessageUsecase(chatter domain.Chatter) OnBasicMessageUsecase {
	return OnBasicMessageUsecase{
		chatter: chatter,
	}
}

func (u OnCommandUsecase) HandleBasicMessage(ctx context.Context, update domain.TGUpdate) {

}
