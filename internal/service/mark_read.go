package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"nox-ai/internal/service/model"
)

func (s *Service) MarkRead(ctx context.Context, data model.WhatsAppStatus) error {

	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}

	header := make(map[string][]string)
	header["Authorization"] = []string{"Bearer " + s.wa.Token}
	header["Content-Type"] = []string{"application/json"}

	_, err = s.http.Request(http.MethodPost, fmt.Sprintf("/%s/%s/messages", s.wa.Version, s.wa.Number), dataByte, header)
	if err != nil {
		return err
	}

	return nil
}
