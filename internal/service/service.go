package service

import (
	"github.com/FaturFawkes/NOX-AI/domain/service"
	"github.com/FaturFawkes/NOX-AI/pkg/config"
	"github.com/go-resty/resty/v2"
	"github.com/sashabaranov/go-openai"
)

type Service struct {
	gpt        *openai.Client
	httpClient *resty.Client
	wa         config.Whatsapp
}

func NewService(gpt *openai.Client, client *resty.Client, wa config.Whatsapp) service.IService {
	return &Service{
		gpt:        gpt,
		httpClient: client,
		wa:         wa,
	}
}
