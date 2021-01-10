package service

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"timelyship.com/accounts/application"
)

type HTTPClient interface {
	SendWithContext(method, url string, body io.Reader) ([]byte, error)
}

type HTTPClientImpl struct {
}

func (h *HTTPClientImpl) SendWithContext(method, url string, body io.Reader) ([]byte, error) {
	// Change NewRequest to NewRequestWithContext and pass context it
	ctx, cancel := context.WithTimeout(context.Background(), application.IntConst.ExternalCallMaxThreshold)
	defer cancel()
	req, newReqErr := http.NewRequestWithContext(ctx, method, url, body)
	if newReqErr != nil {
		return nil, newReqErr
	}
	response, err := http.DefaultClient.Do(req)
	if err == nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(response.Body)
}
