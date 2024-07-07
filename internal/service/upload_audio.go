package service

import (
	"encoding/json"
	"fmt"
	"os"
)

func (s *Service) UploadAudio(path string) (*string, error) {
	// Buka kembali file sementara
	audioFile, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening temp file:", err)
		return nil, err
	}
	defer audioFile.Close()

	buffer := make([]byte, 512)
	_, err = audioFile.Read(buffer)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	resp, err := s.httpClient.R().
		SetAuthToken(s.wa.Token).
		SetFormData(map[string]string{
			"messaging_product": "whatsapp",
		}).
		SetFile("file", path).
		Post(s.wa.Host + fmt.Sprintf("/%s/%s/media", s.wa.Version, s.wa.Number))
	if err != nil {
		return nil, err
	}

	var data struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, err
	}

	return &data.Id, nil
}
