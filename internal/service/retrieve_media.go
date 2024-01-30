package service

import (
	"encoding/json"
	"fmt"
	"github.com/FaturFawkes/NOX-AI/internal/service/model"
)

func (s *Service) RetrieveMedia(audioId string) (string, error) {
	resp, err := s.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(s.wa.Token).
		Get(s.wa.Host + fmt.Sprintf("/%s/%s?phone_number_id=%s", s.wa.Version, audioId, s.wa.Number))
	if err != nil {
		return "nil", err
	}

	var data model.MediaFile

	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return "", err
	}

	return data.URL, nil
}
