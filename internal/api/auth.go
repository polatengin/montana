package api

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/polatengin/montana/internal/config"
)

type Auth struct {
	config *config.ProviderConfig
}

func (client *Auth) GetTokenForScopes(ctx context.Context, scopes []string) (*string, error) {
	tflog.Debug(ctx, fmt.Sprintf("[GetTokenForScope] Getting token for scope: '%s'", strings.Join(scopes, ",")))

	if client.config.Credentials.TestMode {
		token := "test_mode_mock_token_value"
		return &token, nil
	}

	token := ""
	var err error

	switch {
	case client.config.Credentials.IsClientSecretCredentialsProvided():
		token, _, err = client.AuthenticateClientSecret(ctx, scopes)
	case client.config.Credentials.IsCliProvided():
		token, _, err = client.AuthenticateUsingCli(ctx, scopes)
	default:
		return nil, errors.New("no credentials provided")
	}

	if err != nil {
		return nil, err
	}

	return &token, err
}

func (client *Auth) AuthenticateUsingCli(ctx context.Context, scopes []string) (string, time.Time, error) {
	azureCLICredentials, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		return "", time.Time{}, err
	}

	accessToken, err := azureCLICredentials.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: scopes,
	})
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken.Token, accessToken.ExpiresOn, nil
}

func (client *Auth) AuthenticateClientSecret(ctx context.Context, scopes []string) (string, time.Time, error) {
	clientSecretCredential, err := azidentity.NewClientSecretCredential(
		client.config.Credentials.TenantId,
		client.config.Credentials.ClientId,
		client.config.Credentials.ClientSecret, nil)
	if err != nil {
		return "", time.Time{}, err
	}

	accessToken, err := clientSecretCredential.GetToken(ctx, policy.TokenRequestOptions{
		Scopes:   scopes,
		TenantID: client.config.Credentials.TenantId,
	})

	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken.Token, accessToken.ExpiresOn, nil
}
