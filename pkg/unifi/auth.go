package unifi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"path"
)

func Login(cfg Config, username, password string) (tokenCookie *http.Cookie, _ error) {
	data, err := json.Marshal(LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	cfg.Endpoint.Path = path.Join(cfg.Endpoint.Path, "api/auth/login")
	req, err := http.NewRequest(http.MethodPost, cfg.Endpoint.String(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	cookies := response.Cookies()
	for _, ck := range cookies {
		if ck.Name == "TOKEN" {
			tokenCookie = ck
			break
		}
	}
	if tokenCookie == nil {
		return nil, errors.New("Server did not send token cookie in response")
	}
	return
}
