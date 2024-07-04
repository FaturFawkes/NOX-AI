package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
)

func (s *Service) UploadAudio(audio io.ReadCloser) (*string, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	// Menambahkan file audio ke form data
	part, err := writer.CreateFormFile("audio", String(10)+".ogg")
	if err != nil {
		return nil, fmt.Errorf("could not create form file: %v", err)
	}

	// Menyalin audio dari response body ke bagian form
	_, err = io.Copy(part, audio)
	if err != nil {
		return nil, fmt.Errorf("could not copy audio to form: %v", err)
	}

	// Menutup writer untuk menyelesaikan penulisan multipart form data
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close writer: %v", err)
	}

	resp, err := s.httpClient.R().
		SetHeader("Content-Type", "audio/ogg").
		SetAuthToken(s.wa.Token).
		SetBody(payload).
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
