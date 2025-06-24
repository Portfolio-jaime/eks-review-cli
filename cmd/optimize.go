package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// optimizeCmd represents the optimize command
var optimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Recomendaciones para optimizar recursos",
	Long: `Ejecuta análisis sobre el uso de recursos del clúster
y propone ajustes para mejorar el rendimiento y reducir costos.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("optimize called")
	},
}

func init() {
	rootCmd.AddCommand(optimizeCmd)
}
