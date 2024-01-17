package usecase

import (
	"context"
	"github.com/FaturFawkes/NOX-AI/domain/entity"
	"github.com/FaturFawkes/NOX-AI/internal/service/model"
	"go.uber.org/zap"
	"os"
)

func (u *Usecase) HandleAudio(ctx context.Context, user *entity.User, messageId, audioId string) error {

	err := u.service.MarkRead(model.WhatsAppStatus{
		MessagingProduct: "whatsapp",
		Status:           "read",
		MessageID:        messageId,
	})
	if err != nil {
		u.logger.Error("Error mark read message", zap.Error(err))
		panic(err)
	}

	// Get Audio Data
	audioUrl, err := u.service.RetrieveMedia(audioId)
	if err != nil {
		u.logger.Error("Error retrieve audio", zap.Error(err))
		return err
	}

	audioPath, err := u.service.DownloadMedia(audioUrl)
	if err != nil {
		u.logger.Error("Error download audio", zap.Error(err))
		return err
	}

	text, err := u.service.TranscriptionGPT(ctx, audioPath)
	if err != nil {
		u.logger.Error("Error transcription audio", zap.Error(err))
		return err
	}

	// Delete downloaded audio file
	err = os.Remove(audioPath)
	if err != nil {
		u.logger.Error("Error deleting audio file", zap.Error(err))
		return err
	}

	err = u.HandleText(ctx, user, messageId, text)
	if err != nil {
		return err
	}

	return nil
}
