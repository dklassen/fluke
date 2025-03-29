package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var sourceFile string
	var targetFile string

	var rootCmd = &cobra.Command{
		Use:   "fluke",
		Short: "Fluke the ClickHouse migration SQL generator.",
	}

	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Generate migration SQL from two schema files.",
		Run: func(cmd *cobra.Command, args []string) {
			sourceSQLBytes, err := os.ReadFile(sourceFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read source schema file: %v\n", err)
				os.Exit(1)
			}
			targetSQLBytes, err := os.ReadFile(targetFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read target schema file: %v\n", err)
				os.Exit(1)
			}

			sourceSchema := NewSQLParser(string(sourceSQLBytes)).Parse()
			targetSchema := NewSQLParser(string(targetSQLBytes)).Parse()
			migrationSQL := GenerateMigration(sourceSchema, targetSchema)

			fmt.Println(migrationSQL)
		},
	}

	migrateCmd.Flags().StringVarP(&sourceFile, "source", "s", "", "Path to source SQL schema file")
	migrateCmd.Flags().StringVarP(&targetFile, "target", "t", "", "Path to target SQL schema file")
	migrateCmd.MarkFlagRequired("source")
	migrateCmd.MarkFlagRequired("target")

	rootCmd.AddCommand(migrateCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Execute()
}
