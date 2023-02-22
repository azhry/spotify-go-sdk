package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	SpotifyClientX struct {
		baseURL string
		http    *http.Client
	}
)

func New(httpClient *http.Client) *SpotifyClientX {
	cl := &SpotifyClientX{
		http:    httpClient,
		baseURL: "https://api.spotify.com/v1",
	}

	return cl
}

func (cl *SpotifyClientX) Get(ctx context.Context, url string, result interface{}) error {
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	response, err := cl.http.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(result); err != nil {
		return err
	}

	return nil
}

func (cl *SpotifyClientX) CurrentUser(ctx context.Context) error {
	var result interface{}
	if err := cl.Get(ctx, cl.baseURL+"/me", &result); err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
