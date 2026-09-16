package pkg

import (
	"os"
	"path/filepath"

	"github.com/mattn/go-mastodon"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func StoreCredentials(credentials *Credentials) error {
	configFile := viper.ConfigFileUsed()

	if configFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return errors.Wrap(err, "could not determine home directory")
		}

		configDir := filepath.Join(homeDir, ".mastoid")
		if err := os.MkdirAll(configDir, 0700); err != nil {
			return errors.Wrap(err, "could not create config directory")
		}

		configFile = filepath.Join(configDir, "config.yml")
	}

	viper.Set("client_id", credentials.Application.ClientID)
	viper.Set("client_secret", credentials.Application.ClientSecret)
	viper.Set("auth_uri", credentials.Application.AuthURI)
	viper.Set("redirect_uri", credentials.Application.RedirectURI)
	viper.Set("grant_token", credentials.GrantToken)
	viper.Set("access_token", credentials.AccessToken)
	viper.Set("server", credentials.Server)

	if err := viper.WriteConfigAs(configFile); err != nil {
		return errors.Wrap(err, "error writing config")
	}

	if err := os.Chmod(configFile, 0600); err != nil {
		return errors.Wrap(err, "error setting config file permissions")
	}

	return nil
}

func LoadCredentials() (*Credentials, error) {
	clientId := viper.GetString("client_id")
	clientSecret := viper.GetString("client_secret")
	authUri := viper.GetString("auth_uri")
	redirectUri := viper.GetString("redirect_uri")
	grantToken := viper.GetString("grant_token")
	accessToken := viper.GetString("access_token")
	server := viper.GetString("server")

	if clientId == "" || clientSecret == "" {
		return nil, errors.Errorf("no credentials found")
	}

	app := &Credentials{
		Server:      server,
		GrantToken:  grantToken,
		AccessToken: accessToken,
		Application: &mastodon.Application{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			AuthURI:      authUri,
			RedirectURI:  redirectUri,
		},
	}

	return app, nil
}

// Credentials stores the result of an oauth flow.
type Credentials struct {
	Server      string
	GrantToken  string
	Application *mastodon.Application
	AccessToken string
}
