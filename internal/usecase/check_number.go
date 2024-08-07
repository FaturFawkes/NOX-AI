package usecase

import (
	"context"
	"github.com/FaturFawkes/NOX-AI/domain/entity"
	"github.com/FaturFawkes/NOX-AI/internal/service/model"
	"github.com/sashabaranov/go-openai"

	"go.uber.org/zap"
)

func (u *Usecase) CheckNumber(ctx context.Context, data *entity.User) (*entity.User, error) {
	user, err := u.repo.GetUser(data.Number)
	if err != nil {
		user, err = u.repo.InsertUser(data)
		if err != nil {
			u.logger.Error("Error insert user to db", zap.Error(err))
			return nil, err
		}

		err = u.service.SendWA(model.MessageTemplate{
			MessagingProduct: "whatsapp",
			RecipientType:    "individual",
			To:               user.Number,
			Type:             "template",
			Template: model.Template{
				Name: "welcoming",
				Language: model.Language{
					Code: "id",
				},
				Components: []model.Component{
					{
						Type: "header",
						Parameters: []model.Parameter{
							{
								Type: "image",
								Image: model.Image{
									Link: "https://encrypted-tbn1.gstatic.com/images?q=tbn:ANd9GcTx2vzvr55BK2WtWaPrR77WFG-bYYdwaWpPabND_6MwRJUlO7Gl",
								},
							},
						},
					},
				},
			},
		})
		if err != nil {
			u.logger.Error("Error sending greeting message", zap.Error(err))
			return nil, err
		}

	}

	return user, nil
}

func (u *Usecase) InserUserLog(user *entity.User, resGpt *openai.ChatCompletionResponse) error {
	// Add token logger
	res, err := u.repo.GetUserLog(user.ID)
	if err != nil {
		u.logger.Info("No user log for " + user.Number)
		err = u.repo.InsertUserLog(&entity.UserLog{
			UserID:        user.ID,
			TokenRequest:  resGpt.Usage.PromptTokens,
			TokenResponse: resGpt.Usage.CompletionTokens,
			TokenUsage:    resGpt.Usage.TotalTokens,
			TotalRequest:  1,
		})
		if err != nil {
			u.logger.Error("Error insert user log", zap.Error(err))
		}
		return err
	} else {
		res.TokenRequest += resGpt.Usage.PromptTokens
		res.TokenResponse += resGpt.Usage.CompletionTokens
		res.TokenUsage += resGpt.Usage.TotalTokens
		res.TotalRequest++
		err = u.repo.UpdateUserLog(res)
		if err != nil {
			u.logger.Error("Error update log user log", zap.Error(err))
			return err
		}
	}

	return err
}
