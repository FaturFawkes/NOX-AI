package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Service) SendWA(ctx context.Context, number string, data any) error {

	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Println("INI REQUEST ", string(dataByte))

	header := make(map[string][]string)
	header["Authorization"] = []string{"Bearer " + s.wa.Token}
	header["Content-Type"] = []string{"application/json"}

	res, err := s.http.Request(http.MethodPost, fmt.Sprintf("/%s/%s/messages", s.wa.Version, s.wa.Number), dataByte, header)
	if err != nil {
		return err
	}

	fmt.Println("INI RESPONSE ", string(res))
	return nil
}
