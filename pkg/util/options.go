package util

import (
	"os"

	"github.com/spf13/cobra"
)

type Options struct {
	Cwd          string
	DbUser       string
	DbPass       string
	DbHost       string
	DbPort       string
	DbName       string
	DbCollection string
	Debug        bool
}

var options *Options = nil

func envOrDefault(key, defaultValue string) string {
	if envValue, exists := os.LookupEnv(key); exists {
		return envValue
	}

	return defaultValue
}

func InitFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("cwd", "c", ".", "Path working directory, filepaths will be relative to this")
	cmd.Flags().String("dbuser", envOrDefault("DB_USER", "admin"), "Database username")
	cmd.Flags().String("dbpass", envOrDefault("DB_PASS", "admin"), "Database password")
	cmd.Flags().String("dbhost", envOrDefault("DB_HOST", "db"), "Database hostname")
	cmd.Flags().String("dbport", envOrDefault("DB_PORT", "27017"), "Database port")
	cmd.Flags().String("dbname", envOrDefault("DB_NAME", "archidex"), "Database name")
	cmd.Flags().String("dbcollection", envOrDefault("DB_COLLECTION", "entities"), "Database collection")
	cmd.Flags().Bool("debug", false, "Enable debug logging")
}

func InitOptions(cmd *cobra.Command) Options {
	options = &Options{
		Cwd:          cmd.Flag("cwd").Value.String(),
		DbUser:       cmd.Flag("dbuser").Value.String(),
		DbPass:       cmd.Flag("dbpass").Value.String(),
		DbHost:       cmd.Flag("dbhost").Value.String(),
		DbPort:       cmd.Flag("dbport").Value.String(),
		DbName:       cmd.Flag("dbname").Value.String(),
		DbCollection: cmd.Flag("dbcollection").Value.String(),
		Debug:        cmd.Flag("debug").Value.String() == "true",
	}

	return *options
}

func GetOptions() Options {
	return *options
}
