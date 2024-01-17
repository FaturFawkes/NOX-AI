package service

import (
	"github.com/FaturFawkes/NOX-AI/domain/service"
	"github.com/FaturFawkes/NOX-AI/pkg/client"
	"github.com/FaturFawkes/NOX-AI/pkg/config"
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
