package cmd

import (
	"fmt"
	"meta-x/internal"
	"meta-x/lib"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

var mysqlCommand = &cobra.Command{
	Use:   "mysql",
	Short: "use mysql as the database provider",
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

		cfg := mysql.Config{
			User:   dbUsername,
			Passwd: dbPassword,
			DBName: dbName,
			Net:    "tcp",
			Addr:   fmt.Sprintf("%s:%d", dbHost, dbPort),
		}

		port, _ := cmd.Flags().GetInt("port")

		if err := internal.InitDBAndServer(lib.PSQL, cfg.FormatDSN(), port); err != nil {
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
	_ = mysqlCommand.MarkFlagRequired("username")

	rootCmd.AddCommand(mysqlCommand)
}
