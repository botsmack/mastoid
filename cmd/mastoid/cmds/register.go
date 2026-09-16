package cmds

import (
	"context"
	"github.com/go-go-golems/mastoid/cmd/mastoid/pkg"
	"github.com/mattn/go-mastodon"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register the app with the Mastodon instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		clientName, _ := cmd.Flags().GetString("client-name")
		redirectURIs, _ := cmd.Flags().GetString("redirect-uris")
		scopes, _ := cmd.Flags().GetString("scopes")
		website, _ := cmd.Flags().GetString("website")
		server, _ := cmd.Flags().GetString("server")

		ctx := context.Background()
		appConfig := &mastodon.AppConfig{
			ClientName:   clientName,
			RedirectURIs: redirectURIs,
			Scopes:       scopes,
			Website:      website,
			Server:       server,
		}

		log.Info().Msg("Registering app...")
		app, err := mastodon.RegisterApp(ctx, appConfig)
		if err != nil {
			log.Error().Err(err).Msgf("Error registering app")
			return errors.Wrap(err, "Error registering app")
		}

		credentials := &pkg.Credentials{
			Server:      server,
			Application: app,
		}

		log.Info().Msg("App registration successful!\n")
		log.Debug().Str("ClientID", app.ClientID).Msgf("Client ID")
		log.Debug().Str("AuthURI", app.AuthURI).Msgf("Auth URI")
		log.Debug().Str("RedirectURI", app.RedirectURI).Msgf("Redirect URI")

		log.Info().Msg("Authorizing app...")

		err = pkg.Authorize(ctx, credentials)
		if err != nil {
			log.Error().Err(err).Msgf("Error authorizing app")
			return errors.Wrap(err, "Error authorizing app")
		}

		err = pkg.StoreCredentials(credentials)
		if err != nil {
			return errors.Wrap(err, "Error storing credentials")
		}

		log.Info().Msg("App authorization successful!")

		return nil
	},
}
