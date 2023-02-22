package auth

import (
	"context"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
)

const (
	AuthURL  = "https://accounts.spotify.com/authorize"
	TokenURL = "https://accounts.spotify.com/api/token"
)

type (
	Authenticator struct {
		config *oauth2.Config
	}
)

func New() *Authenticator {
	config := &oauth2.Config{
		ClientID:     "96eec0bec03d4a60b843a8c9c4e82137",
		ClientSecret: "25a23e96b6b849a4be0031e274eacc3e",
		Scopes:       []string{"user-read-private", "user-read-email", "user-top-read", "user-read-currently-playing", "user-read-playback-state", "user-modify-playback-state", "app-remote-control", "streaming", "playlist-modify-private", "playlist-modify-public"},
		RedirectURL:  "http://localhost:3000/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  AuthURL,
			TokenURL: TokenURL,
		},
	}

	authenticator := &Authenticator{
		config: config,
	}

	return authenticator
}

func (authenticator Authenticator) GetAuthURL(state string) string {
	return authenticator.config.AuthCodeURL(state)
}

func (authenticator Authenticator) GetToken(ctx context.Context, state string, r *http.Request) (*oauth2.Token, error) {
	values := r.URL.Query()
	if err := values.Get("error"); err != "" {
		return nil, errors.New(err)
	}

	code := values.Get("code")

	return authenticator.config.Exchange(ctx, code)
}

func (authenticator Authenticator) Client(ctx context.Context, token *oauth2.Token) *http.Client {
	return authenticator.config.Client(ctx, token)
}
