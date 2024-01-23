package service

func (s *Service) TranscribeYoutube(url string, lang string) (int, string, error) {
	resp, err := s.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"videoUrl": url,
			"langCode": lang,
		}).
		Post("https://tactiq-apps-prod.tactiq.io/transcript")
	if err != nil {
		return 500, "nil", err
	}

	return resp.StatusCode(), string(resp.Body()), nil
}
