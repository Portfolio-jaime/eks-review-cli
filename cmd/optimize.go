package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// optimizeCmd represents the optimize command
var optimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Sugiere mejoras de rendimiento",
	Long: `Identifica recursos sin uso o poco aprovechados y revisa la
configuración de autoescalado para proponer ajustes que optimicen el
consumo del clúster.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("optimize called")
	},
}

func init() {
	rootCmd.AddCommand(optimizeCmd)
}
