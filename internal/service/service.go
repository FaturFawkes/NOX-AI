package service

import (
	"nox-ai/domain/service"
	"nox-ai/pkg/client"
	"nox-ai/pkg/config"

	"github.com/sashabaranov/go-openai"
)

type Service struct {
	gpt  *openai.Client
	http client.ClientHttp
	wa   config.Whatsapp
}

func NewService(gpt *openai.Client, client client.ClientHttp, wa config.Whatsapp) service.IService {
	return &Service{
		gpt:  gpt,
		http: client,
		wa:   wa,
	}
}
