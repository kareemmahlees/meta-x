package cmd

import (
	"meta-x/internal"

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
		if err = internal.InitDBAndServer("sqlite3", filePath, port); err != nil {
			return err
		}
		return nil
	},
}

func init() {

	sqlite3Command.Flags().StringP("file", "f", "", "database file path")
	sqlite3Command.MarkFlagRequired("file")

	rootCmd.AddCommand(sqlite3Command)
}
