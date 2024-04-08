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
	Short: "use sqlite as the database provider",
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}
		sqliteConfig := utils.NewSQLiteConfig(filePath)

		conn, err := internal.InitDBConn(lib.SQLITE3, sqliteConfig)
		if err != nil {
			return err
		}
		provider := db.NewSQLiteProvider(conn)

		server := internal.NewServer(provider, port, make(chan bool))
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
