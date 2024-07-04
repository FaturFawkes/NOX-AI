package service

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

func (s *Service) AudioGPT(ctx context.Context, text string) (*openai.RawResponse, error) {
	res, err := s.gpt.CreateSpeech(ctx, openai.CreateSpeechRequest{
		Model:          openai.TTSModel1,
		Input:          text,
		Voice:          openai.VoiceNova,
		ResponseFormat: openai.SpeechResponseFormatOpus,
		Speed:          1,
	})

	if err != nil {
		return nil, err
	}

	return &res, nil
}
