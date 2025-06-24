package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Herramientas de monitoreo del clúster",
	Long: `Agrupa subcomandos orientados a la observación y seguimiento del estado
del clúster de Kubernetes.
Permite revisar métricas y eventos para detectar problemas de forma temprana.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("monitor called")
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
