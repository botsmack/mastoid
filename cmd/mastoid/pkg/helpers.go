package pkg

import (
	"context"
	"github.com/mattn/go-mastodon"
	"github.com/rs/zerolog/log"
	"regexp"
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

// ExtractID matches [0-9]+ in a string and returns the matched value
func ExtractID(status string) string {
	reg := regexp.MustCompile(`[0-9]+`)
	return reg.FindString(status)
}
