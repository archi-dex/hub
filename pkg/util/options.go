package util

import (
	"os"

	"github.com/spf13/cobra"
)

type Options struct {
	Cwd    string
	DbUser string
	DbPass string
	DbHost string
	DbPort string
	DbName string
	Debug  bool
}

var (
	options *Options
)

func envOrDefault(key, defaultValue string) string {
	if envValue, exists := os.LookupEnv(key); exists {
		return envValue
	}

	return defaultValue
}

func InitFlags(cmd *cobra.Command) {
	if cmd == nil {
		return
	}

	cmd.Flags().String("cwd", ".", "Path working directory, filepaths will be relative to this")
	cmd.Flags().String("dbuser", envOrDefault("DB_USER", "admin"), "Database username")
	cmd.Flags().String("dbpass", envOrDefault("DB_PASS", "admin"), "Database password")
	cmd.Flags().String("dbhost", envOrDefault("DB_HOST", "db"), "Database hostname")
	cmd.Flags().String("dbport", envOrDefault("DB_PORT", "27017"), "Database port")
	cmd.Flags().String("dbname", envOrDefault("DB_NAME", "archidex"), "Database name")
	cmd.Flags().Bool("debug", false, "Enable debug logging")
}

func InitOptions(cmd *cobra.Command) Options {
	if cmd == nil {
		options = &Options{
			Cwd:    ".",
			DbUser: envOrDefault("DB_USER", "admin"),
			DbPass: envOrDefault("DB_PASS", "admin"),
			DbHost: envOrDefault("DB_HOST", "db"),
			DbPort: envOrDefault("DB_PORT", "27017"),
			DbName: envOrDefault("DB_NAME", "archidex"),
			Debug:  true,
		}
	} else {
		options = &Options{
			Cwd:    cmd.Flag("cwd").Value.String(),
			DbUser: cmd.Flag("dbuser").Value.String(),
			DbPass: cmd.Flag("dbpass").Value.String(),
			DbHost: cmd.Flag("dbhost").Value.String(),
			DbPort: cmd.Flag("dbport").Value.String(),
			DbName: cmd.Flag("dbname").Value.String(),
			Debug:  cmd.Flag("debug").Value.String() == "true",
		}
	}

	return *options
}

func GetOptions() Options {
	return *options
}
