package service

import (
	"context"
	"github.com/FaturFawkes/NOX-AI/domain/entity"
	"github.com/FaturFawkes/NOX-AI/internal/service/model"
	"github.com/sashabaranov/go-openai"
	"io"
)

type IService interface {
	TextGPT(ctx context.Context, model string, message []openai.ChatCompletionMessage) (*openai.ChatCompletionResponse, error)
	ImageGPT(ctx context.Context, prompt string) (*openai.ImageResponse, error)
	AudioGPT(ctx context.Context, text string) (*openai.RawResponse, error)
	SendWA(data any) error
	MarkRead(data model.WhatsAppStatus) error
	RetrieveMedia(audioId string) (string, error)
	UploadAudio(audio io.ReadCloser) (*string, error)
	DownloadMedia(link string, mediaType entity.TypeMedia) (string, error)
	TranscriptionGPT(ctx context.Context, path string) (string, error)
	TranscribeYoutube(url string, lang string) (int, string, error)
}
