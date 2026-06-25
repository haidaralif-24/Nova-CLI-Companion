package cmd

import (
	"fmt"
	"nova/core/memory"
	"os"

	"github.com/spf13/cobra"
)

var logListFlag bool

var logCmd = &cobra.Command{
	Use: "log",
	Run: func(cmd *cobra.Command, args []string) {
		if logListFlag {
			listLoggedCommands()
			return
		}

		if len(args) == 0 {
			fmt.Println("no command provided")
			return
		}

		command := args[0]
		cwd, _ := os.Getwd()

		_, err := memory.DB.Exec(
			"INSERT INTO commands (command, cwd, exit_code) VALUES (?, ?, ?)",
			command, cwd, 0,
		)
		if err != nil {
			fmt.Println("failed to log command:", err)
			return
		}

		fmt.Println("logged:", command)
	},
}

func listLoggedCommands() {
	rows, err := memory.DB.Query(
		"SELECT id, command, cwd, exit_code, timestamp FROM commands ORDER BY id DESC LIMIT 20",
	)
	if err != nil {
		fmt.Println("failed to read commands:", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, exitCode int
		var command, cwd, timestamp string
		if err := rows.Scan(&id, &command, &cwd, &exitCode, &timestamp); err != nil {
			fmt.Println("failed to read row:", err)
			return
		}
		fmt.Printf("[%d] %s  (exit=%d)  %s  in %s\n", id, command, exitCode, timestamp, cwd)
		count++
	}

	if count == 0 {
		fmt.Println("no commands logged yet")
	}
}

func init() {
	logCmd.Flags().BoolVarP(&logListFlag, "list", "l", false, "list recent logged commands")
	rootCmd.AddCommand(logCmd)
}
