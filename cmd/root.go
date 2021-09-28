package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "uff",
	Short: "Unfriend Facebook Friends smartly :)",
	Long:  "Scrap your friend list from facebook and unfriend them as you like",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println(`
         ________
  __  __/ __/ __/
 / / / / /_/ /_  
/ /_/ / __/ __/  
\__,_/_/ /_/        
				`)

		fmt.Println("Unfriend Facebook Friends - smartly")

		fmt.Println()
		fmt.Println()

		fmt.Println("Make sure you have .env file configured")
		fmt.Println("Use 'uuf scrap' to start")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
