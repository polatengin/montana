package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/polatengin/montana/internal/config"
)

type ProviderClient struct {
	Config *config.ProviderConfig
	Api    *ApiClient
}

type ApiHttpResponse struct {
	Response    *http.Response
	BodyAsBytes []byte
}

func (client *ApiClient) GetConfig() *config.ProviderConfig {
	return client.Config
}

type ApiClient struct {
	Config   *config.ProviderConfig
	BaseAuth *Auth
}

func NewApiClientBase(config *config.ProviderConfig, baseAuth *Auth) *ApiClient {
	return &ApiClient{
		Config:   config,
		BaseAuth: baseAuth,
	}
}

func TryGetScopeFromURL(url string) (string, error) {
	switch {
	case strings.LastIndex(url, "api.bap.microsoft.com") != -1,
		strings.LastIndex(url, "api.powerapps.com") != -1:

		return "https://service.powerapps.com/.default", nil
	case strings.LastIndex(url, "api.powerplatform.com") != -1:

		return "https://api.powerplatform.com/.default", nil
	case strings.LastIndex(url, ".com/") != -1:

		scope := strings.SplitAfterN(url, ".com/", 2)[0]
		scope = scope + ".default"
		return scope, nil
	default:
		return "", errors.New("Unable to determine scope from url: '" + url + "'. Please provide your own scope.")
	}
}

func (client *ApiClient) Execute(ctx context.Context, method string, url string, headers http.Header, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error) {
	scope, err := TryGetScopeFromURL(url)
	if err != nil {
		return nil, err
	}

	token, err := client.BaseAuth.GetTokenForScopes(ctx, []string{scope})
	if err != nil {
		return nil, err
	}

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
	apiResponse, err := client.doRequest(token, request, headers)
	if err != nil {
		return nil, err
	}

	isStatusCodeValid := false
	for _, statusCode := range acceptableStatusCodes {
		if apiResponse.Response.StatusCode == statusCode {
			isStatusCodeValid = true
			break
		}
	}
	if !isStatusCodeValid {
		return nil, fmt.Errorf("expected status code: %d, recieved: %d", acceptableStatusCodes, apiResponse.Response.StatusCode)
	}
	if responseObj != nil {
		err = apiResponse.MarshallTo(responseObj)
		if err != nil {
			return nil, err
		}
	}
	return apiResponse, nil
}

func (client *ApiClient) doRequest(token *string, request *http.Request, headers http.Header) (*ApiHttpResponse, error) {
	apiHttpResponse := &ApiHttpResponse{}
	if headers != nil {
		request.Header = headers
	}

	if token == nil || *token == "" {
		return nil, errors.New("token is empty")
	}

	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	httpClient := http.DefaultClient

	if request.Header["Authorization"] == nil {
		request.Header.Set("Authorization", "Bearer "+*token)
	}

	request.Header.Set("User-Agent", "terraform-provider-power-platform")

	response, err := httpClient.Do(request)
	apiHttpResponse.Response = response
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	apiHttpResponse.BodyAsBytes = body
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if len(body) != 0 {
			return apiHttpResponse, fmt.Errorf("status: %d, message: %s", response.StatusCode, string(body))
		} else {
			return nil, fmt.Errorf("status: %d", response.StatusCode)
		}
	}
	return apiHttpResponse, nil
}

func (apiResponse *ApiHttpResponse) MarshallTo(obj interface{}) error {
	err := json.NewDecoder(bytes.NewReader(apiResponse.BodyAsBytes)).Decode(&obj)
	if err != nil {
		return err
	}
	return nil
}
