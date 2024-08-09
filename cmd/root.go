/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/kznagamori/go_har2xlsx/lib"
	"github.com/spf13/cobra"
)

var (
	inputFile  string
	outputFile string
	verbose    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go_har2xlsx",
	Short: ".harファイルをエクセルファイルに変換する",
	Long:  `go_har2xlsxはHTTP Archive (.har)形式の記録をエクセルファイル(.xlsx)形式に変換します。`,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
		} else {
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
		}
		err := lib.ExecuteHar2xlsx(inputFile, outputFile)
		if err != nil {
			slog.Error("%v", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true // completionコマンド無効化
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init 関数は、パッケージがインポートされる際、プログラムが開始される前に自動的に実行されます。
// main 関数よりも先に実行されます。
func init() {
	// Define `--file` (`-f`) flag
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file")

	// Define `--out` (`-o`) flag
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file")

	// Mark flags as required if needed
	rootCmd.MarkFlagRequired("inputFile")
	rootCmd.MarkFlagRequired("outputFile")

	// Define global `--verbose` flag
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable	verbose	output")
}
