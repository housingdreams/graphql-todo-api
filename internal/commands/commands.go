package commands

import (
	"fmt"
	"net/http"
	"strings"

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
		// Search config in home directory with name ".cobra" (no extension).
		viper.AddConfigPath("./conf")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/todo")
		viper.SetConfigName("todo")
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
	viper.SetDefault("server.hostname", ":5555")
	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.name", "todo")
	viper.SetDefault("database.user", "minh")
	viper.SetDefault("database.password", "anhyeuem98")

	viper.SetDefault("queue.broker", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("queue.store", "memcache://localhost:11211")

	rootCmd.SetVersionTemplate(versionTemplate)
	rootCmd.AddCommand(newWebCmd(), newMigrateCommand(), newTokenCmd())
	rootCmd.Execute()
}
