package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run shopify customer grpc command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("grpc called")
	},
}

func init() {
	serveCmd.AddCommand(grpcCmd)
}
