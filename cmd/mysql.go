package cmd

import (
	"github.com/kareemmahlees/meta-x/internal"
	"github.com/kareemmahlees/meta-x/internal/db"
	"github.com/kareemmahlees/meta-x/lib"
	"github.com/kareemmahlees/meta-x/utils"
	"github.com/spf13/cobra"
)

var mysqlCommand = &cobra.Command{
	Use:   "mysql",
	Short: "use mysql as the database provider",
	RunE: func(cmd *cobra.Command, args []string) error {
		var mysqlConfig *utils.MySQLConfig

		port, _ := cmd.Flags().GetInt("port")

		connUrl, _ := cmd.Flags().GetString("url")
		if connUrl != "" {
			mysqlConfig = utils.NewMySQLConfig(&connUrl, nil)
		} else {
			dbUsername, _ := cmd.Flags().GetString("username")
			dbHost, _ := cmd.Flags().GetString("host")
			dbPort, _ := cmd.Flags().GetInt("dbPort")
			dbName, _ := cmd.Flags().GetString("db")
			dbPassword, _ := cmd.Flags().GetString("password")

			mysqlConfig = utils.NewMySQLConfig(nil, &utils.MySQLConnectionParams{
				DBUsername: dbUsername,
				DBPassword: dbPassword,
				DBHost:     dbHost,
				DBPort:     dbPort,
				DBName:     dbName,
			})
		}

		conn, err := internal.InitDBConn(lib.MYSQL, mysqlConfig)
		if err != nil {
			return err
		}
		provider := db.NewMySQLProvider(conn)
		server := internal.NewServer(provider, port, make(chan<- bool))

		if err := server.Serve(); err != nil {
			return err
		}

		return nil
	},
}

func init() {

	mysqlCommand.Flags().String("username", "root", "db username")
	mysqlCommand.Flags().String("password", "", "db password")
	mysqlCommand.Flags().String("host", "localhost", "db host")
	mysqlCommand.Flags().Int("dbPort", 3306, "db port")
	mysqlCommand.Flags().String("db", "mysql", "db name")
	mysqlCommand.Flags().String("url", "", "connection url/string")

	mysqlCommand.MarkFlagsMutuallyExclusive("password", "url")

	rootCmd.AddCommand(mysqlCommand)
}
