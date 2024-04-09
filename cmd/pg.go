package cmd

import (
	"github.com/kareemmahlees/meta-x/internal"
	"github.com/kareemmahlees/meta-x/internal/db"
	"github.com/kareemmahlees/meta-x/lib"
	"github.com/kareemmahlees/meta-x/utils"
	"github.com/spf13/cobra"
)

var pgCommand = &cobra.Command{
	Use:   "pg",
	Short: "use postgres as the database provider",
	RunE: func(cmd *cobra.Command, args []string) error {
		var pgConfig utils.PgConfig

		port, _ := cmd.Flags().GetInt("port")

		connUrl, _ := cmd.Flags().GetString("url")
		if connUrl != "" {
			pgConfig = utils.PgConfig{ConnUrl: &connUrl, ConnParams: nil}
		} else {
			dbUsername, _ := cmd.Flags().GetString("username")
			dbHost, _ := cmd.Flags().GetString("host")
			dbPort, _ := cmd.Flags().GetInt("dbPort")
			dbName, _ := cmd.Flags().GetString("db")
			dbSslMode, _ := cmd.Flags().GetString("sslmode")
			dbPassword, _ := cmd.Flags().GetString("password")

			pgConfig = utils.PgConfig{ConnUrl: nil, ConnParams: &utils.PgConnectionParams{
				DbUsername: dbUsername,
				DbPassword: dbPassword,
				DbHost:     dbHost,
				DbPort:     dbPort,
				DbName:     dbName,
				DbSslMode:  dbSslMode,
			}}
		}
		conn, err := internal.InitDBConn(lib.PSQL, pgConfig)
		if err != nil {
			return err
		}
		provider := db.NewPgProvider(conn)
		server := internal.NewServer(provider, port, make(chan<- bool))

		if err := server.Serve(); err != nil {
			return err
		}

		return nil
	},
}

func init() {

	pgCommand.Flags().String("username", "postgres", "db username")
	pgCommand.Flags().String("password", "", "db password")
	pgCommand.Flags().String("host", "localhost", "db host")
	pgCommand.Flags().Int("dbPort", 5432, "db port")
	pgCommand.Flags().String("db", "postgres", "db name")
	pgCommand.Flags().String("url", "", "connection url/string")
	pgCommand.Flags().String("sslmode", "disable", "db sslmode")
	pgCommand.MarkFlagsMutuallyExclusive("password", "url")

	rootCmd.AddCommand(pgCommand)
}
