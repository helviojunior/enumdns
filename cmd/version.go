package cmd

import (
    "fmt"

    "github.com/helviojunior/enumdns/internal/ascii"
    "github.com/helviojunior/enumdns/internal/version"
    "github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Get the enumdns version",
    Long:  ascii.LogoHelp(`Get the enumdns version.`),
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println(ascii.Logo())

        fmt.Println("Author: Helvio Junior (m4v3r1ck)")
        fmt.Println("Source: https://github.com/helviojunior/enumdns")
        fmt.Printf("Version: %s\nGit hash: %s\nBuild env: %s\nBuild time: %s\n\n",
            version.Version, version.GitHash, version.GoBuildEnv, version.GoBuildTime)
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}