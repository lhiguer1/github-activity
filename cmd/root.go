/*
Copyright © 2024 Leonel Higuera <lhiguer1@asu.edu>
*/
package cmd

import (
	"net/http"
	"os"
	"time"

	. "github.com/lhiguer1/github-activity/github_activity"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "github-activity",
	Annotations: map[string]string{
		cobra.CommandDisplayNameAnnotation: "github-activity username",
	},
	Short: "Fetch recent activity of a specified GitHub user.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := &http.Client{Timeout: 10 * time.Second}
		service := NewGitHubService(client)
		username := args[0]
		PrintRecentActivity(username, service)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
