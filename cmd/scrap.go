package cmd

import (
	"os"
	"uff/scrap"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scrapCommand)
}

var scrapCommand = &cobra.Command{
	Use:   "scrap",
	Short: "Scrap friend list from facebook",
	Run: func(cmd *cobra.Command, args []string) {
		scrap.FetchFriends(os.Getenv("url"))
	},
}
