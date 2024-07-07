package service

import (
	"fmt"
	"github.com/FaturFawkes/NOX-AI/domain/entity"
	"github.com/FaturFawkes/NOX-AI/pkg/utils"
)

func (s *Service) DownloadMedia(link string, mediaType entity.TypeMedia) (string, error) {
	var path string

	if mediaType == entity.TypeAudio {
		path = fmt.Sprintf("%s.ogg", utils.String(10))
	} else {
		path = fmt.Sprintf("%s.png", utils.String(10))
	}
	_, err := s.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(s.wa.Token).
		SetOutput(path).
		Get(link)
	if err != nil {
		return "", err
	}

	return path, nil
}
