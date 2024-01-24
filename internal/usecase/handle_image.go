package usecase

import (
	"context"
	"fmt"
	"github.com/FaturFawkes/NOX-AI/domain/entity"
	"github.com/FaturFawkes/NOX-AI/internal/delivery/request"
)

func (u *Usecase) HandleImage(ctx context.Context, user *entity.User, image request.Image) error {
	fmt.Println("INI IMAGE ", image)
	return nil
}
