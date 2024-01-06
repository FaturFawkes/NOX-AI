package config

import "github.com/sashabaranov/go-openai"

func NewGPT(apikeyGPT string) *openai.Client {
	return openai.NewClient(apikeyGPT)
}