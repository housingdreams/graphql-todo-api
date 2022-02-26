package commands

// import (
// 	"fmt"
// 	"time"

// 	"github.com/jmoiron/sqlx"
// 	"github.com/leminhson2398/todo-api/internal/db"
// 	log "github.com/sirupsen/logrus"
// 	"github.com/spf13/cobra"
// 	"github.com/spf13/viper"
// )

// func newResetPasswordCmd() *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "reset-password",
// 		Short: "Resets password of the specified user",
// 		Long:  "If the user forgets it's password you can reset it with this command",
// 		Args:  cobra.ExactArgs(2),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			connection := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable",
// 				viper.GetString("database.user"),
// 				viper.GetString("database.password"),
// 				viper.GetString("database.host"),
// 				viper.GetString("database.name"),
// 			)
// 			var database *sqlx.DB
// 			var err error
// 			var retryDuration time.Duration
// 			maxTryNumber := 4
// 			for i := 0; i < maxTryNumber; i++ {
// 				database, err = sqlx.Connect("postgres", connection)
// 				if err == nil {
// 					break
// 				}
// 				retryDuration = time.Duration(i*2) * time.Second
// 				log.WithFields(log.Fields{
// 					"retryNumber":   i,
// 					"retryDuration": retryDuration.Unix(),
// 				}).WithError(err).Error("issue connecting the database")
// 				if i != maxTryNumber-1 {
// 					time.Sleep(retryDuration)
// 				}
// 			}
// 			database.SetMaxOpenConns(25)
// 			database.SetMaxIdleConns(25)
// 			database.SetConnMaxLifetime(5 * time.Minute)
// 			repo := *db.NewRepository(database)

// 			username := args[0]
// 			password := args[1]

// 			user, err := repo.GetUserAccountByID()
// 		},
// 	}
// }
