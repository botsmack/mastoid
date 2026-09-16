package cmds

import (
	"context"
	"github.com/go-go-golems/mastoid/cmd/mastoid/pkg"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var AuthorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "Authorize the app with the Mastodon instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		credentials, err := pkg.LoadCredentials()
		if err != nil {
			return errors.Wrap(err, "Error loading credentials")
		}

		log.Info().Msg("Authorizing app...")
		err = pkg.Authorize(context.Background(), credentials)
		if err != nil {
			log.Error().Err(err).Msgf("Error authorizing app")
			return errors.Wrap(err, "Error authorizing app")
		}

		log.Info().Msg("App authorization successful!\n")

		err = pkg.StoreCredentials(credentials)
		if err != nil {
			return errors.Wrap(err, "Error storing credentials")
		}

		return nil
	},
}
