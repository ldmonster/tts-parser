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
		Short: "Tool for parsing tabletop simulator modules",
		Long:  `Tool for parsing tabletop simulator modules`,
	}

	var text = &cobra.Command{
		Use:   "text",
		Short: "Tool for parsing tabletop simulator modules",
		Long:  `Tool for parsing tabletop simulator modules`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("bebbe")
		},
	}

	rootCmd.AddCommand(text)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
