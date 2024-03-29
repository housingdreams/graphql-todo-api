package commands

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/leminhson2398/todo-api/internal/auth"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newTokenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "token [userid]",
		Short: "Create a long lived JWT token for dev purpose",
		Long:  "Create a long lived JWT token for dev purpose",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			secret := viper.GetString("server.secret")
			if strings.TrimSpace(secret) == "" {
				return errors.New("server.secret must be set")
			}

			// Use userid provided
			id := args[0]

			// make one up
			if id == "" {
				id = uuid.New().String()
				fmt.Println("generated", id)
			}

			token, err := auth.NewAccessTokenCustomExpiration(id, time.Hour*24, []byte(secret))
			if err != nil {
				log.WithError(err).Error("issue while creating access token")
				return err
			}
			fmt.Println(token)
			return nil
		},
	}
}
