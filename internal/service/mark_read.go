package service

import (
	"fmt"
	"github.com/FaturFawkes/NOX-AI/internal/service/model"
)

func (s *Service) MarkRead(data model.WhatsAppStatus) error {

	_, err := s.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(s.wa.Token).
		SetBody(`{"status": "read"}`).
		Put(s.wa.Host + fmt.Sprintf("/%s/messages/%s", s.wa.Version, data.MessageID))
	if err != nil {
		return err
	}

	return nil
}
