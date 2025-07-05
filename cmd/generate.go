/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Pitchouneee/koomos/internal/diagram"
	"github.com/Pitchouneee/koomos/internal/parser"
	"github.com/Pitchouneee/koomos/internal/resolver"
	"github.com/spf13/cobra"
)

var input string
var output string

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a diagram from Kubernetes/ArgoCD YAML files",
	Long:  `Scans a directory recursively to parse Kubernetes, ArgoCD, Kustomize or Helm manifests and generates a Mermaid diagram.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Scanning YAML files from:", input)

		resources, err := parser.ParseDirectory(input)
		if err != nil {
			fmt.Println("Error parsing YAML files:", err)
			os.Exit(1)
		}

		if len(resources) == 0 {
			fmt.Println("No YAML resources found.")
			os.Exit(0)
		}

		edges := resolver.ResolveRelations(resources)
		diagramText := diagram.GenerateMermaid(resources, edges)

		err = os.WriteFile(output, []byte(diagramText), 0644)
		if err != nil {
			fmt.Println("Failed to write output file:", err)
			os.Exit(1)
		}

		fmt.Printf("Mermaid diagram generated with %d resources: %s\n", len(resources), output)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	generateCmd.Flags().StringVarP(&input, "input", "i", ".", "Dossier contenant les fichiers YAML")
	generateCmd.Flags().StringVarP(&output, "output", "o", "diagram.md", "Fichier de sortie (Markdown avec Mermaid)")
}
