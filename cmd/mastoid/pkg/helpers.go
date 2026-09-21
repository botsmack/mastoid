package pkg

import (
	"context"
	"github.com/mattn/go-mastodon"
	"github.com/rs/zerolog/log"
	"net/url"
	"path"
	"strings"
)

func CreateClient(credentials *Credentials) (*mastodon.Client, error) {
	config := &mastodon.Config{
		Server:       credentials.Server,
		ClientID:     credentials.Application.ClientID,
		ClientSecret: credentials.Application.ClientSecret,
		AccessToken:  credentials.AccessToken,
	}
	client := mastodon.NewClient(config)

	return client, nil
}

func CreateClientAndAuthenticate(ctx context.Context, credentials *Credentials) (*mastodon.Client, error) {
	client, err := CreateClient(credentials)
	if err != nil {
		return nil, err
	}

	if client.Config.AccessToken == "" {
		log.Debug().Msg("Authenticating with app")

		err = client.AuthenticateApp(ctx)
		if err != nil {
			return nil, err
		}

		log.Debug().Msg("App authenticated")
	} else {
		log.Debug().Msg("Already has access token")
	}

	_, err = client.VerifyAppCredentials(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// ExtractID returns a Mastodon status ID from either a numeric ID or a status URL.
func ExtractID(status string) string {
	status = strings.TrimSpace(status)
	if status == "" {
		return ""
	}

	if isNumericID(status) {
		return status
	}

	parsed, err := url.Parse(status)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}

	statusID := path.Base(strings.TrimRight(parsed.Path, "/"))
	if !isNumericID(statusID) {
		return ""
	}

	return statusID
}

func isNumericID(value string) bool {
	if value == "" {
		return false
	}

	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}
