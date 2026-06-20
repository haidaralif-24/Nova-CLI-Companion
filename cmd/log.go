package cmd

import (
	"fmt"
	"nova/core/memory"
	"os"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use: "log",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("no command provided")
			return
		}

		command := args[0]
		cwd, _ := os.Getwd()

		memory.DB.Exec(
			"INSERT INTO commands (command, cwd, exit_code) VALUES (?, ?, ?)",
			command, cwd, 0,
		)

		fmt.Println("logged:", command)
	},
}
