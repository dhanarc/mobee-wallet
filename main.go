package main

import (
	"github.com/dhanarc/mobee-wallet/cmd"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	http := &cobra.Command{
		Use:   "http",
		Short: "Run HTTP Server",
		Run:   cmd.HTTP,
	}

	rootCmd := &cobra.Command{Use: "app"}

	rootCmd.AddCommand(http)
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
