package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/caarlos0/solarman-exporter/config"
	"golang.org/x/oauth2"
)

type Client struct {
	c   *http.Client
	cfg config.Config
}

func New(cfg config.Config) (*Client, error) {
	auth, err := newAccessToken(cfg)
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	if err := json.Unmarshal(auth, &token); err != nil {
		return nil, fmt.Errorf("could not auth: %w", err)
	}

	oauthConfg := oauth2.Config{
		ClientID:     cfg.AppID,
		ClientSecret: cfg.AppSecret,
	}

	c := oauthConfg.Client(context.Background(), &token)
	return &Client{
		c:   c,
		cfg: cfg,
	}, nil
}

func newAccessToken(cfg config.Config) ([]byte, error) {
	data := fmt.Sprintf(
		`{"appSecret":%q,"email":%q,"password":%q}`,
		cfg.AppSecret,
		cfg.Username,
		cfg.Password,
	)

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"https://globalapi.solarmanpv.com/account/v1.0/token?appId=%s&language=en&=",
			cfg.AppID,
		),
		strings.NewReader(data),
	)
	if err != nil {
		return nil, fmt.Errorf("could not auth: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not auth: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not auth: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not auth: %w", err)
	}

	return bts, nil
}
