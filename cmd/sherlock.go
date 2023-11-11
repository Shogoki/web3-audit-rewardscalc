/*-
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
    "log"
    "shogoki/audit-rewardscalc/parser"
    "os"
    "io"
	"github.com/spf13/cobra"
)

// sherlockCmd represents the sherlock command
var sherlockCmd = &cobra.Command{
	Use:   "sherlock",
	Short: "Get Rewards for Sherlock Contest",
    Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
        repoName := args[0]
        contest, err := parser.GetContestDetails(repoName)
        if(err != nil) {
            log.Fatal("Failed to get contest Details", err)
        }

        stdIn, err := io.ReadAll(os.Stdin)
     if(err != nil) {
        panic(err)
    }
    contest.Issues, err = parser.GetContestIssues(stdIn)
    if(err != nil) {
        log.Fatal("Failed to get Issues")
    }
    totalShares := contest.GetTotalShares()
    fmt.Printf("Contest %s\nTOTAL Issues %2d | PrizePool: %6d USDC | Total Shares: %6.3f\n",repoName, len(contest.Issues), contest.PrizePool, contest.GetTotalShares())
    fmt.Println()
    var sanityCheck float64
    for _, issue := range contest.Issues {
        shares := issue.GetShares()
        percentage := shares / totalShares 
        reward := percentage * float64(contest.PrizePool)
        sanityCheck += reward * float64((len(issue.Duplicates)+1))
        fmt.Printf("Issue:  %5d \t| Duplicates: %5d \t| %s -\t%s\n",issue.Number,len(issue.Duplicates), issue.Severity,issue.Title)
        fmt.Printf("Shares: %5.3f\t| Percentage: %5.2f%%\t| Payout per Watson: %9.3f USDC\n",shares, percentage*100,reward)
        fmt.Println()
        }
    if(contest.PrizePool != uint32(sanityCheck)) {
        log.Fatal("ATTENTION: Total Payouts do not match Prizepool!!")
    }

	},
}


func init() {
	rootCmd.AddCommand(sherlockCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sherlockCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sherlockCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
