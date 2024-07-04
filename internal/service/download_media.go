package service

import (
	"fmt"
	"github.com/FaturFawkes/NOX-AI/domain/entity"
	"math/rand"
	"time"
)

func (s *Service) DownloadMedia(link string, mediaType entity.TypeMedia) (string, error) {
	var path string

	if mediaType == entity.TypeAudio {
		path = fmt.Sprintf("%s.ogg", String(10))
	} else {
		path = fmt.Sprintf("%s.png", String(10))
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

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}
