package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

func newPasswordCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "password",
		Short: "Generate bcrypt password of the specified string",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			password := args[0]

			hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			fmt.Println(string(hash))
			return nil
		},
	}
}
