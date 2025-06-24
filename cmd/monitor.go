package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Herramientas para inspeccionar el clúster",
	Long: `Agrupa subcomandos orientados a observar el estado del clúster.
Permite listar recursos, revisar eventos, acceder a logs y obtener
resúmenes generales de la salud de Kubernetes.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("monitor called")
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
