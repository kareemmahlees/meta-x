package cmd

import (
	"fmt"
	"meta-x/internal"
	"meta-x/lib"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

var mysqlCommand = &cobra.Command{
	Use:   "mysql",
	Short: "use mysql as the database provider",
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

			dbPassword, _ := cmd.Flags().GetString("password")
			if dbPassword == "" {
				fmt.Println("Enter password: ")
				fmt.Scanln(&dbPassword)
			}

			conf := mysql.Config{
				User:   dbUsername,
				Passwd: dbPassword,
				DBName: dbName,
				Net:    "tcp",
				Addr:   fmt.Sprintf("%s:%d", dbHost, dbPort),
			}
			cfg = conf.FormatDSN()
		}

		port, _ := cmd.Flags().GetInt("port")

		app := fiber.New(fiber.Config{DisableStartupMessage: true})

		if err := internal.InitDBAndServer(app, lib.MYSQL, cfg, port, make(chan bool, 1)); err != nil {
			return err
		}
		return nil
	},
}

func init() {

	mysqlCommand.Flags().String("username", "", "db username")
	mysqlCommand.Flags().String("password", "", "db password")
	mysqlCommand.Flags().String("host", "localhost", "db host")
	mysqlCommand.Flags().Int("dbPort", 3306, "db port")
	mysqlCommand.Flags().String("db", "mysql", "db name")
	mysqlCommand.Flags().String("url", "", "connection url/string")

	mysqlCommand.MarkFlagsMutuallyExclusive("username", "url")

	rootCmd.AddCommand(mysqlCommand)
}
