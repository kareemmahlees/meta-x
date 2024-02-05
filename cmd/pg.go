package cmd

import (
	"fmt"
	"meta-x/internal"
	"meta-x/lib"

	"github.com/lib/pq"

	"github.com/spf13/cobra"
)

var pgCommand = &cobra.Command{
	Use:   "pg",
	Short: "use postgres as the database provider",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbPassword, _ := cmd.Flags().GetString("password")
		if dbPassword == "" {
			fmt.Println("Enter password: ")
			fmt.Scanln(&dbPassword)
		}

		dbUsername, _ := cmd.Flags().GetString("username")
		dbHost, _ := cmd.Flags().GetString("host")
		dbPort, _ := cmd.Flags().GetInt("dbPort")
		dbName, _ := cmd.Flags().GetString("db")
		dbSslMode, _ := cmd.Flags().GetString("sslmode")

		cfg, _ := pq.ParseURL(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", dbUsername, dbPassword, dbHost, dbPort, dbName, dbSslMode))

		port, _ := cmd.Flags().GetInt("port")

		if err := internal.InitDBAndServer(lib.PSQL, cfg, port); err != nil {
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
	pgCommand.Flags().String("sslmode", "disable", "db sslmode")
	_ = pgCommand.MarkFlagRequired("username")

	rootCmd.AddCommand(pgCommand)
}
