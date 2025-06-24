package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// diagnoseCmd represents the diagnose command
var diagnoseCmd = &cobra.Command{
	Use:   "diagnose",
	Short: "Utilidades de diagnóstico",
	Long: `Revisa el estado de componentes del clúster
y ayuda a localizar fallas o inconsistencias en los recursos.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("diagnose called")
	},
}

func init() {
	rootCmd.AddCommand(diagnoseCmd)
}
