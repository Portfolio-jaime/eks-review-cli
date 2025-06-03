// cmd/get.go
package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd representa el comando 'monitor get'
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Muestra uno o más recursos por tipo",
	Long: `Muestra información detallada sobre uno o más tipos de recursos de Kubernetes
(pods, services, daemonsets, etc.).

Similar a 'kubectl get'.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Puedes poner lógica aquí que aplique a todos los subcomandos de 'get'
		// si es necesario en el futuro.
	},
}

func init() {
	monitorCmd.AddCommand(getCmd) // Añade 'get' como subcomando de 'monitor'

	// Añadir todos los nuevos subcomandos a 'getCmd'
	getCmd.AddCommand(podsGetCmd)
	getCmd.AddCommand(servicesGetCmd)
	getCmd.AddCommand(daemonsetsGetCmd)
	getCmd.AddCommand(jobsGetCmd)
	getCmd.AddCommand(cronjobsGetCmd)
	getCmd.AddCommand(namespacesGetCmd)
	getCmd.AddCommand(serviceaccountsGetCmd)
}
