package cluster

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	_contentTypeHeader = "Content-Type"
	_contentTypeJson   = "application/json"
)

type BaseClient interface {
	Request(ctx context.Context, method, url string, body any, params map[string]string) (*http.Request, *bytes.Buffer, error)
	Do(request *http.Request, body *bytes.Buffer, response any) error
}

type Impl struct {
	client  http.Client
	baseUrl string
}

func NewBaseClient(baseUrl string) BaseClient {
	return &Impl{
		client:  http.Client{},
		baseUrl: baseUrl,
	}
}

func (c *Impl) Request(ctx context.Context, method, uri string, body any, params map[string]string) (*http.Request, *bytes.Buffer, error) {
	requestUrl, err := c.prepareUrl(uri, params)
	if err != nil {
		return nil, nil, err
	}

	requestBody := &bytes.Buffer{}
	if body != nil {
		requestBody, err = toBytes(body)
		if err != nil {
			return nil, nil, err
		}
	}

	httpRequest, err := http.NewRequestWithContext(ctx, method, requestUrl, requestBody)
	if err != nil {
		return nil, nil, err
	}

	httpRequest.Header.Add(_contentTypeHeader, _contentTypeJson)

	return httpRequest, requestBody, nil
}

func (c *Impl) Do(request *http.Request, body *bytes.Buffer, response any) error {
	httpResponse, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("[%s] не удалось выполнить запрос, %w", request.URL.String(), err)
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("[%s] не удалось выполнить запрос, код ответа: %d", request.URL.String(), httpResponse.StatusCode)
	}

	if response != nil {
		if err = json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
			return fmt.Errorf("[%s] не удалось сериализовать ответ, %w", request.URL.String(), err)
		}
	}

	return nil
}

func (c *Impl) prepareUrl(uri string, params map[string]string) (string, error) {
	requestUrl := fmt.Sprintf("%s/%s", c.baseUrl, uri)
	if len(params) == 0 {
		return requestUrl, nil
	}

	parsedUrl, err := url.Parse(requestUrl)
	if err != nil {
		return "", err
	}

	query := parsedUrl.Query()

	for k, v := range params {
		if len(v) == 0 {
			continue
		}

		query.Add(k, v)
	}

	parsedUrl.RawQuery = query.Encode()

	return parsedUrl.String(), nil
}

func toBytes(value any) (*bytes.Buffer, error) {
	bytesBody, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(bytesBody), err
}
