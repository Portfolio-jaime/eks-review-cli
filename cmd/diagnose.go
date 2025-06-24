package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// diagnoseCmd represents the diagnose command
var diagnoseCmd = &cobra.Command{
	Use:   "diagnose",
	Short: "Ejecuta chequeos de diagnóstico",
	Long: `Realiza pruebas básicas en Pods, Services e Ingresses para
detectar configuraciones erróneas o estados anómalos.
Actualmente solo imprime un mensaje de prueba.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("diagnose called")
	},
}

func init() {
	rootCmd.AddCommand(diagnoseCmd)
}
