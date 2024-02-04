package cmd

import (
	"meta-x/internal"
	"meta-x/lib"

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
		if err = internal.InitDBAndServer(lib.SQLITE3, filePath, port); err != nil {
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
