package commands

import (
	"fmt"
	"net/http"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const mainDescription = `Todo is an open source project`

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionTemplate = fmt.Sprintf(`Version: %s
Commit: %s
Built: %s`, version, commit, date+"\n")

var cfgFile string

var rootCmd = &cobra.Command{
	Use:     "todo",
	Long:    mainDescription,
	Version: version,
}

var migration http.FileSystem

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
	migration = http.Dir("./migrations")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigFile(".env")
	}

	viper.SetEnvPrefix("TODO")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Execute the root cobra command
func Execute() {
	// viper.SetDefault("server.hostname", ":8000")
	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.name", "docker")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "docker")

	rootCmd.SetVersionTemplate(versionTemplate)
	rootCmd.AddCommand(newWebCmd(), newMigrateCommand(), newTokenCmd(), newPasswordCmd())
	rootCmd.Execute()
}
