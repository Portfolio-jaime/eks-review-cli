package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Verbose permitirá la salida detallada (logs de DEBUG)
var Verbose bool // Exportada para ser accesible por otros archivos del paquete cmd

var rootCmd = &cobra.Command{
	Use:   "eks-review",
	Short: "Una herramienta CLI para revisar clústeres de EKS.",
	Long: `eks-review es una herramienta de línea de comandos (CLI) escrita en Go,
diseñada para simplificar la revisión y el diagnóstico de recursos en clústeres de Kubernetes,
con un enfoque especial en Amazon EKS.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Habilitar salida detallada (logs de DEBUG)")
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle") // Quita esto si no se usa
}
