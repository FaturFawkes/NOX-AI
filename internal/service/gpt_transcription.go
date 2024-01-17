package service

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

func (s *Service) TranscriptionGPT(ctx context.Context, path string) (string, error) {
	res, err := s.gpt.CreateTranscription(ctx, openai.AudioRequest{
		Model:       openai.Whisper1,
		FilePath:    path,
		Temperature: 0,
	})
	if err != nil {
		return "", err
	}

	return res.Text, nil
}
