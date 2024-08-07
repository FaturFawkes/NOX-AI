package service

import (
	"encoding/json"
	"fmt"
)

func (s *Service) SendWA(data any) error {

	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := s.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(s.wa.Token).
		SetBody(dataByte).
		Post(s.wa.Host + fmt.Sprintf("/%s/%s/messages", s.wa.Version, s.wa.Number))
	if err != nil {
		return err
	}

	fmt.Println("INI RESPONSE WA ", string(resp.Body()))

	return nil
}
