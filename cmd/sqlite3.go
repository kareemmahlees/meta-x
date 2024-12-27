package cmd

import (
	"github.com/kareemmahlees/meta-x/internal"
	"github.com/kareemmahlees/meta-x/internal/db"
	"github.com/kareemmahlees/meta-x/lib"
	"github.com/kareemmahlees/meta-x/utils"

	"github.com/spf13/cobra"
)

var sqlite3Command = &cobra.Command{
	Use:   "sqlite3",
	Short: "Use SQLite as the database provider",
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}
		// Not handling the error because port always has a default value, `5522`.
		port, _ := cmd.Flags().GetInt("port")

		sqliteConfig := utils.NewSQLiteConfig(filePath)

		conn, err := db.InitDBConn(lib.SQLITE3, sqliteConfig)
		if err != nil {
			return err
		}
		provider := db.NewSQLiteProvider(conn)

		server := internal.NewServer(provider, port, make(chan bool, 1))
		if err = server.Serve(); err != nil {
			return err
		}
		return nil
	},
}

func init() {

	sqlite3Command.Flags().StringP("file", "f", "", "database file path")
	_ = sqlite3Command.MarkFlagRequired("file")

	rootCmd.AddCommand(sqlite3Command)
}
