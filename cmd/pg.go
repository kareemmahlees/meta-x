package cmd

import (
	"fmt"
	"meta-x/internal"
	"meta-x/lib"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"

	"github.com/spf13/cobra"
)

var pgCommand = &cobra.Command{
	Use:   "pg",
	Short: "use postgres as the database provider",
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg string

		connUrl, _ := cmd.Flags().GetString("url")
		if connUrl != "" {
			cfg = connUrl
		} else {
			dbUsername, _ := cmd.Flags().GetString("username")
			dbHost, _ := cmd.Flags().GetString("host")
			dbPort, _ := cmd.Flags().GetInt("dbPort")
			dbName, _ := cmd.Flags().GetString("db")
			dbSslMode, _ := cmd.Flags().GetString("sslmode")

			dbPassword, _ := cmd.Flags().GetString("password")
			if dbPassword == "" {
				fmt.Println("Enter password: ")
				fmt.Scanln(&dbPassword)
			}

			cfg, _ = pq.ParseURL(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", dbUsername, dbPassword, dbHost, dbPort, dbName, dbSslMode))
		}

		port, _ := cmd.Flags().GetInt("port")

		app := fiber.New(fiber.Config{DisableStartupMessage: true})

		if err := internal.InitDBAndServer(app, lib.PSQL, cfg, port, make(chan bool, 1)); err != nil {
			return err
		}
		return nil
	},
}

func init() {

	pgCommand.Flags().String("username", "", "db username")
	pgCommand.Flags().String("password", "", "db password")
	pgCommand.Flags().String("host", "localhost", "db host")
	pgCommand.Flags().Int("dbPort", 5432, "db port")
	pgCommand.Flags().String("db", "postgres", "db name")
	pgCommand.Flags().String("url", "", "connection url/string")
	pgCommand.Flags().String("sslmode", "disable", "db sslmode")
	pgCommand.MarkFlagsMutuallyExclusive("username", "url")

	rootCmd.AddCommand(pgCommand)
}
