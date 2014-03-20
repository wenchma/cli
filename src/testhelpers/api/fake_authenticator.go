package api

import (
	"cf/configuration"
	"cf/errors"
)

type FakeAuthenticationRepository struct {
	Config           configuration.ReadWriter
	AuthenticateArgs struct {
		Credentials []map[string]string
	}
	GetLoginPromptsReturns struct {
		Error   error
		Prompts map[string]configuration.AuthPrompt
	}

	AuthError    bool
	AccessToken  string
	RefreshToken string
}

func (auth *FakeAuthenticationRepository) Authenticate(credentials map[string]string) (apiErr error) {
	auth.AuthenticateArgs.Credentials = append(auth.AuthenticateArgs.Credentials, copyMap(credentials))

	if auth.AuthError {
		apiErr = errors.New("Error authenticating.")
		return
	}

	if auth.AccessToken == "" {
		auth.AccessToken = "BEARER some_access_token"
	}

	auth.Config.SetAccessToken(auth.AccessToken)
	auth.Config.SetRefreshToken(auth.RefreshToken)

	return
}

func (auth *FakeAuthenticationRepository) RefreshAuthToken() (updatedToken string, apiErr error) {
	return
}

func (auth *FakeAuthenticationRepository) GetLoginPromptsAndSaveUAAServerURL() (prompts map[string]configuration.AuthPrompt, apiErr error) {
	prompts = auth.GetLoginPromptsReturns.Prompts
	apiErr = auth.GetLoginPromptsReturns.Error
	return
}

func copyMap(input map[string]string) map[string]string {
	output := map[string]string{}
	for key, val := range input {
		output[key] = val
	}
	return output
}
