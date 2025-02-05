/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "tts-parser",
		Short: "Tool for parsing and managing Tabletop Simulator modules",
		Long: `A CLI tool for parsing Tabletop Simulator module files (.json), downloading assets,
creating backups, and auditing downloaded files.`,
	}

	var downloadCmd = &cobra.Command{
		Use:   "download [module_path]",
		Short: "Download module assets",
		Long:  `Download all assets referenced in a Tabletop Simulator module file`,
		Run: func(cmd *cobra.Command, args []string) {
			start()
		},
	}

	var backupCmd = &cobra.Command{
		Use:   "backup [module_path]",
		Short: "Backup module files",
		Long:  `Create a backup of module files and downloaded assets`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Backing up module files...")
		},
	}

	var auditCmd = &cobra.Command{
		Use:   "audit [module_path]",
		Short: "Audit module files",
		Long:  `Check integrity of downloaded module assets`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Auditing module files...")
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringP("temp-dir", "t", "tmp/", "Temporary download directory")
	rootCmd.PersistentFlags().DurationP("timeout", "o", 0, "Download timeout duration (e.g. 30s, 1m)")
	rootCmd.PersistentFlags().BoolP("overwrite", "w", false, "Overwrite existing files when downloading")

	// Backup command flags
	backupCmd.Flags().StringP("output", "o", "backups/", "Backup output directory")

	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(auditCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
