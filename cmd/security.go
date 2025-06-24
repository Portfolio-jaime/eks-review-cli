package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// securityCmd represents the security command
var securityCmd = &cobra.Command{
	Use:   "security",
	Short: "Audita configuraciones de seguridad",
	Long: `Evalúa Network Policies, reglas RBAC, imágenes de contenedores
y Secrets para detectar posibles vulnerabilidades o malas prácticas.
Actualmente es un comando placeholder sin funcionalidad completa.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("security called")
	},
}

func init() {
	rootCmd.AddCommand(securityCmd)
}
