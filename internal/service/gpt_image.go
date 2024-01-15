package service

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

func (s *Service) ImageGPT(ctx context.Context, prompt string) (*openai.ImageResponse, error) {
	resp, err := s.gpt.CreateImage(ctx, openai.ImageRequest{
		Prompt: prompt,
		Model: openai.CreateImageModelDallE3,
	})
	if err != nil {
		return nil, err
	}
	
	return &resp, nil
}