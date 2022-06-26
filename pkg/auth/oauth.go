package auth

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
)

type Auth struct {
	conf        clientcredentials.Config
	tokenSource oauth2.TokenSource
}

func New(clientId, clientSecret string) *Auth {
	conf := clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{"READSYSTEM", "WRITESYSTEM"},
		TokenURL:     "https://api.myuplink.com/oauth/token",
	}
	tokenSource := conf.TokenSource(context.Background())
	return &Auth{
		conf:        conf,
		tokenSource: tokenSource,
	}
}

func (s *Auth) Token() string {
	t, err := s.tokenSource.Token()
	if err != nil {
		// TODO: log me
		return ""
	}
	return t.AccessToken
}

func (s *Auth) Intercept(_ context.Context, req *http.Request) error {
	token, err := s.tokenSource.Token()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	return nil
}
