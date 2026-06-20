package cmd

import (
	"fmt"
	"os"
	"time"

	"nova/context/core"

	"github.com/spf13/cobra"
)

var greetCmd = &cobra.Command{
	Use: "greet",
	Run: func(cmd *cobra.Command, args []string) {
		hour := time.Now().Hour()

		var timeOfDay string
		switch {
		case hour < 12:
			timeOfDay = "Morning"
		case hour < 18:
			timeOfDay = "Afternoon"
		default:
			timeOfDay = "Evening"
		}

		cwd, _ := os.Getwd()
		project := core.DetectProject(cwd)

		fmt.Println("Good", timeOfDay)
		if project != "" {
			fmt.Println("Detected project:", project)
		} else {
			fmt.Println("No project detected")
		}
	},
}

func init() {
	rootCmd.AddCommand(greetCmd)
}
