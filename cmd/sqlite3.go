package cmd

import (
	"meta-x/internal"
	"meta-x/lib"

	"github.com/gofiber/fiber/v2"
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
		app := fiber.New(fiber.Config{DisableStartupMessage: true})
		if err = internal.InitDBAndServer(app, lib.SQLITE3, filePath, port, make(chan bool)); err != nil {
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
