package service

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

func (s *Service) TextGPT(ctx context.Context, message []openai.ChatCompletionMessage) (string, error) {
	resp, err := s.gpt.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo1106,
		Messages: message,
	})
	if err != nil {
		fmt.Println("[DEBUG AI] error: ", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
