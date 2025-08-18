package cmd

import (
	"fmt"

	"github.com/bob-reis/enumdns/internal/ascii"
	"github.com/bob-reis/enumdns/internal/version"
	"github.com/spf13/cobra"
)

var releaseOnly = false
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the enumdns version",
	Long:  ascii.LogoHelp(`Get the enumdns version.`),
	Run: func(cmd *cobra.Command, args []string) {
		if releaseOnly {
			fmt.Printf("%s\n",
				version.Version)
		} else {
			fmt.Println(ascii.Logo())

			fmt.Println("Author: Helvio Junior (m4v3r1ck)")
			fmt.Println("Source: https://github.com/bob-reis/enumdns")
			fmt.Printf("Version: %s\nGit hash: %s\nBuild env: %s\nBuild time: %s\n\n",
				version.Version, version.GitHash, version.GoBuildEnv, version.GoBuildTime)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.PersistentFlags().BoolVarP(&releaseOnly, "--release", "r", false, "Show release only")
}
