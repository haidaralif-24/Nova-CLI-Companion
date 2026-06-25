package cmd

import (
	"fmt"

	"nova/core/memory"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nova",
	Short: "Nova — CLI Companion",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return memory.Init()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
