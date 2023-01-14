package sdk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type (
	SpotifyClient struct {
		ApiUrl       string `json:"api_url"`
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		CallbackUri  string `json:"callback_uri"`
		Scopes       string `json:"scopes"`
		StateKey     string `json:"state_key"`
	}
)

func (s *SpotifyClient) Initialize(clientId, clientSecret, callbackUri, scopes string) {
	s.ApiUrl = "https://api.spotify.com"
	s.StateKey = "spotify_auth_state"

	s.ClientId = clientId
	s.ClientSecret = clientSecret
	s.CallbackUri = callbackUri
	s.Scopes = scopes

	envClientId := os.Getenv("SPOTIFY_CLIENT_ID")
	if envClientId != "" {
		s.ClientId = envClientId
	}
	envClientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if envClientSecret != "" {
		s.ClientSecret = envClientSecret
	}
	envCallbackUri := os.Getenv("SPOTIFY_CALLBACK_URI")
	if envCallbackUri != "" {
		s.CallbackUri = envCallbackUri
	}
	envScopes := os.Getenv("SPOTIFY_SCOPES")
	if envScopes != "" {
		s.Scopes = envScopes
	}
}

func (s *SpotifyClient) GetAuthRedirect(state string) string {
	return "https://accounts.spotify.com/authorize?client_id=" + s.ClientId + "&response_type=code&redirect_uri=" + s.CallbackUri + "&scope=" + s.Scopes + "&state=" + state
}

func (s *SpotifyClient) GetClientCredentials(clientId, clientSecret string) (*http.Response, error) {
	param := url.Values{}
	param.Add("grant_type", "client_credentials")

	var payload = bytes.NewBufferString(param.Encode())
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", payload)
	if err != nil {
		return nil, err
	}
	authorization := "Basic " + base64.StdEncoding.EncodeToString([]byte(clientId+":"+clientSecret))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *SpotifyClient) GetAuthToken(code, state string) (*http.Response, error) {
	param := url.Values{}
	param.Add("grant_type", "authorization_code")
	param.Add("code", code)
	param.Add("redirect_uri", s.CallbackUri)

	var payload = bytes.NewBufferString(param.Encode())
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", payload)
	if err != nil {
		return nil, err
	}
	authorization := "Basic " + base64.StdEncoding.EncodeToString([]byte(s.ClientId+":"+s.ClientSecret))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *SpotifyClient) GetArtists(token, ids string) (*http.Response, error) {
	request, err := http.NewRequest("GET", s.ApiUrl+"/v1/artists?ids="+ids, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", token)
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *SpotifyClient) GetFeaturedPlaylists(token string) (*http.Response, error) {
	request, err := http.NewRequest("GET", s.ApiUrl+"/v1/browse/featured-playlists", nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", token)
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *SpotifyClient) CreatePlaylist(token, userId string, data []byte) (*http.Response, error) {
	jsonData := map[string]interface{}{}
	json.Unmarshal(data, &jsonData)
	delete(jsonData, "user_id")
	delete(jsonData, "track_uris")
	jsonByte, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", s.ApiUrl+"/v1/users/"+userId+"/playlists", bytes.NewBuffer(jsonByte))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", token)
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *SpotifyClient) AddItemToPlaylist(token, playlistId string, tracks []string) (*http.Response, error) {
	data := struct {
		Uris []string `json:"uris"`
	}{
		Uris: tracks,
	}
	jsonByte, _ := json.Marshal(data)

	request, err := http.NewRequest("POST", s.ApiUrl+"/v1/playlists/"+playlistId+"/tracks", bytes.NewBuffer(jsonByte))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", token)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *SpotifyClient) GetTracksRecommendations(token string, seedArtists []string, seedGenres []string, seedTracks []string) (*http.Response, error) {
	request, err := http.NewRequest("GET", s.ApiUrl+"/v1/recommendations", nil)
	if err != nil {
		return nil, err
	}

	q := request.URL.Query()
	if len(seedArtists) > 0 {
		q.Add("seed_artists", strings.Join(seedArtists, ","))
	}
	if len(seedGenres) > 0 {
		q.Add("seed_genres", strings.Join(seedGenres, ","))
	}
	if len(seedTracks) > 0 {
		q.Add("seed_tracks", strings.Join(seedTracks, ","))
	}

	request.URL.RawQuery = q.Encode()

	request.Header.Add("Authorization", token)
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *SpotifyClient) SearchItem(token, query string, category []string) (*http.Response, error) {
	request, err := http.NewRequest("GET", s.ApiUrl+"/v1/search", nil)
	if err != nil {
		return nil, err
	}

	q := request.URL.Query()
	q.Add("q", query)
	if len(category) == 0 {
		category = []string{"track"}
	}
	q.Add("type", strings.Join(category, ","))

	request.URL.RawQuery = q.Encode()

	request.Header.Add("Authorization", token)
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *SpotifyClient) GetMe(token string) (*http.Response, error) {
	request, err := http.NewRequest("GET", s.ApiUrl+"/v1/me", nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", token)
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *SpotifyClient) GetAlbums(token, ids string) (*http.Response, error) {
	request, err := http.NewRequest("GET", s.ApiUrl+"/v1/albums?ids="+ids, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", token)
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *SpotifyClient) GetTrack(token, id string) (*http.Response, error) {
	request, err := http.NewRequest("GET", s.ApiUrl+"/v1/tracks/"+id, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", token)
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
