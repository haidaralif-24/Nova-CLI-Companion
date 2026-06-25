package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"nova/context/core"
	"nova/llm"

	"github.com/spf13/cobra"
)

var greetCmd = &cobra.Command{
	Use: "greet",
	Run: func(cmd *cobra.Command, args []string) {
		hour := time.Now().Hour()

		var timeOfDay string
		switch {
		case hour < 12:
			timeOfDay = "morning"
		case hour < 18:
			timeOfDay = "afternoon"
		default:
			timeOfDay = "evening"
		}

		cwd, _ := os.Getwd()
		info := core.DetectProjectInfo(cwd)

		if tryAIGreet(timeOfDay, info) {
			return
		}
		templateGreet(timeOfDay, info)
	},
}

// tryAIGreet attempts an AI-generated greeting; falls back if unconfigured.
func tryAIGreet(timeOfDay string, info core.ProjectInfo) bool {
	cfg, err := llm.LoadConfig()
	if err != nil || !cfg.IsConfigured() {
		return false
	}

	provider := llm.FromConfig(*cfg)
	if provider == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	greeting, err := provider.Greet(ctx, llm.GreetRequest{
		TimeOfDay:   timeOfDay,
		ProjectName: info.Name,
		ProjectType: info.Type,
		Cwd:         info.Path,
	})
	if err != nil || strings.TrimSpace(greeting) == "" {
		return false
	}

	fmt.Println(greeting)
	return true
}

// templateGreet prints a static greeting offline.
func templateGreet(timeOfDay string, info core.ProjectInfo) {
	fmt.Println("Good", strings.Title(timeOfDay))
	if info.Type != "" {
		fmt.Printf("%s project detected.\n", info.Type)
		fmt.Println("Ready to continue?")
	} else {
		fmt.Println("No project detected.")
	}
}

func init() {
	rootCmd.AddCommand(greetCmd)
}
