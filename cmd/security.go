package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// securityCmd represents the security command
var securityCmd = &cobra.Command{
	Use:   "security",
	Short: "Comandos para verificar seguridad",
	Long: `Ejecuta revisiones de políticas, permisos y configuraciones
que afectan la seguridad del clúster de Kubernetes.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("security called")
	},
}

func init() {
	rootCmd.AddCommand(securityCmd)
}
