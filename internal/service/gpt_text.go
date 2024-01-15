package service

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

func (s *Service) TextGPT(ctx context.Context, model string, message []openai.ChatCompletionMessage) (*openai.ChatCompletionResponse, error) {
	resp, err := s.gpt.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    model,
		Messages: message,
	})
	if err != nil {
		fmt.Println("[DEBUG AI] error: ", err)
		return nil, err
	}

	return &resp, nil
}
