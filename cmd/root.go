package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "meta-x",
	Short: "RESTfull & Graphql API for your database",
}

func Execute() {

	rootCmd.PersistentFlags().IntP("port", "p", 5522, "port to serve on")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
