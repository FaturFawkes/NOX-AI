package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientHttp interface {
	Request(method string, url string, request any, headers map[string][]string) ([]byte, error)
}

type clientHttp struct {
	ctx    context.Context
	log    *zap.Logger
	host   string
	client HttpClient
}

func NewHttpClient(host string, ctx context.Context, log *zap.Logger, timeout time.Duration) ClientHttp {
	hc := &clientHttp{
		ctx:    ctx,
		log:    log,
		host:   host,
		client: nil,
	}
	proxyUrl := os.Getenv("PROXY_URL")
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	if proxyUrl != "" {
		proxy, err := url.Parse(proxyUrl)
		if err != nil {
			log.Fatal("cant parse url proxy", zap.Error(err))
			os.Exit(1)
		}
		customTransport.Proxy = http.ProxyURL(proxy)
	}
	hc.client = &http.Client{Transport: customTransport, Timeout: timeout}
	return hc
}

func (h *clientHttp) Request(method string, url string, request any, headers map[string][]string) ([]byte, error) {
	requestUrl := h.host + url
	var body io.Reader
	h.log.Info(fmt.Sprintf("making new %s request host:%s , url: %s", method, h.host, url), zap.Any("payload", request))
	body = nil
	if val, ok := request.([]byte); ok {
		body = bytes.NewReader(val)
	}
	httpHeader := h.initiateHttpHeaders(headers)
	httpRequest, err := http.NewRequestWithContext(
		h.ctx,
		method,
		strings.TrimSpace(requestUrl),
		body,
	)
	if err != nil {
		return nil, err
	}
	httpRequest.Header = httpHeader
	res, err := h.client.Do(httpRequest)
	if err != nil {
		h.log.Error("error call http client do", zap.Error(err))
		return nil, err
	}
	rByte, err := h.manageResponseHttp(res, err)
	h.log.Info(fmt.Sprintf("response body from %s", requestUrl), zap.Any("response", string(rByte)))
	return rByte, err
}

func (h *clientHttp) manageResponseHttp(response *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, errors.New(ClassifyNetworkError(err))
	}

	// read response body
	res, err := io.ReadAll(response.Body)
	// leverage defer stack to defer closing of response body read operation
	// this will defer until this function is ready to return

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.log.Error("error cant close body", zap.Error(err))
		}
	}(response.Body)

	if err != nil {
		return res, errors.New("invalid response body")
	}
	// check response status
	if response.StatusCode == 403 {
		return res, errors.New("not authorized")
	}

	if response.StatusCode >= 400 && response.StatusCode < 500 {
		return res, fmt.Errorf("error with code %d", response.StatusCode)
	}
	if response.StatusCode >= 500 {
		return res, errors.New("internal server error")
	}
	return res, nil

}

func (h *clientHttp) initiateHttpHeaders(headers map[string][]string) http.Header {
	httpHeader := make(http.Header)
	if headers != nil {
		return headers
	}
	httpHeader.Set("Content-Type", "application/json")
	httpHeader.Set("Accept", "application/json")
	return httpHeader
}
