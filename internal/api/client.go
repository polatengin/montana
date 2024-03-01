package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ApiClient struct {
}

type ApiHttpResponse struct {
	Response    *http.Response
	BodyAsBytes []byte
}

func (client *ApiClient) Execute(ctx context.Context, method string, url string, headers http.Header, body interface{}, acceptableStatusCodes []int) (*ApiHttpResponse, error) {
	var bodyBuffer io.Reader = nil
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyBuffer = bytes.NewBuffer(bodyBytes)
	}

	request, err := http.NewRequestWithContext(ctx, method, url, bodyBuffer)
	if err != nil {
		return nil, err
	}
	apiHttpResponse := &ApiHttpResponse{}
	if headers != nil {
		request.Header = headers
	}

	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	httpClient := http.DefaultClient

	request.Header.Set("User-Agent", "terraform-provider-power-platform")

	response, err := httpClient.Do(request)
	apiHttpResponse.Response = response
	if err != nil {
		return nil, err
	}

	_body, err := io.ReadAll(response.Body)
	apiHttpResponse.BodyAsBytes = _body
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if len(_body) != 0 {
			return apiHttpResponse, fmt.Errorf("status: %d, message: %s", response.StatusCode, string(_body))
		} else {
			return nil, fmt.Errorf("status: %d", response.StatusCode)
		}
	}

	if err != nil {
		return nil, err
	}

	isStatusCodeValid := false
	for _, statusCode := range acceptableStatusCodes {
		if apiHttpResponse.Response.StatusCode == statusCode {
			isStatusCodeValid = true
			break
		}
	}
	if !isStatusCodeValid {
		return nil, fmt.Errorf("expected status code: %d, recieved: %d", acceptableStatusCodes, apiHttpResponse.Response.StatusCode)
	}
	return apiHttpResponse, nil
}
