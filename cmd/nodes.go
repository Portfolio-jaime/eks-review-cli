package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	// "path/filepath" // No necesario aquí

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/client-go/kubernetes" // Se obtiene de GetKubeClients
	// "k8s.io/client-go/tools/clientcmd" // Se usa en GetKubeClients
	// "k8s.io/client-go/util/homedir" // Se usa en GetKubeClients
	// metrics "k8s.io/metrics/pkg/client/clientset/versioned" // Se obtiene de GetKubeClients
)

// nodesCmd representa el comando 'monitor nodes'
var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Muestra información detallada sobre los nodos del clúster.",
	Long: `El comando nodes recupera y muestra información detallada sobre
los nodos del clúster de Kubernetes, incluyendo su estado, roles y uso de recursos.
El uso de recursos requiere un servidor de métricas instalado en el clúster.`,
	Run: func(cmd *cobra.Command, args []string) {
		clients := GetKubeClients() // 1. Obtener clientes

		// Ya no necesitamos crear metricsClientset explícitamente aquí,
		// GetKubeClients() lo maneja y clients.Metrics será nil si falla.

		fmt.Fprintln(os.Stdout, "Recuperando información de nodos de Kubernetes...")

		// Usar clients.Core
		nodes, err := clients.Core.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listando nodos: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stdout, "\n--- Nodos ---")
		if len(nodes.Items) == 0 {
			fmt.Fprintln(os.Stdout, "No se encontraron nodos.")
		} else {
			headers := []string{"NOMBRE", "ESTADO", "ROLES", "VERSIÓN", "CPU_ALLOC", "CPU_USO", "MEM_ALLOC", "MEM_USO", "EDAD"}
			rows := make([][]string, 0, len(nodes.Items)) // Inicializar rows

			if Verbose { // 2. Usar Verbose
				fmt.Fprintf(os.Stdout, "DEBUG: Nodos encontrados: %d\n", len(nodes.Items))
			}

			for _, node := range nodes.Items {
				status := getNodeStatus(node)
				roles := getNodeRoles(node)
				kubeletVersion := node.Status.NodeInfo.KubeletVersion
				age := metav1.Now().Sub(node.CreationTimestamp.Time).Truncate(time.Second).String() // Usar Truncate

				cpuAlloc := node.Status.Allocatable[corev1.ResourceCPU]
				memAlloc := node.Status.Allocatable[corev1.ResourceMemory]
				cpuAllocStr := cpuAlloc.String()
				memAllocStr := memAlloc.String()

				cpuUsage := "N/A"
				memUsage := "N/A"

				if clients.Metrics != nil { // Usar clients.Metrics
					nodeMetrics, err := clients.Metrics.MetricsV1beta1().NodeMetricses().Get(context.TODO(), node.Name, metav1.GetOptions{})
					if err == nil {
						cpuUsed := nodeMetrics.Usage[corev1.ResourceCPU]
						memUsed := nodeMetrics.Usage[corev1.ResourceMemory]

						if cpuAlloc.MilliValue() > 0 {
							cpuUsage = fmt.Sprintf("%.2f%%", float64(cpuUsed.MilliValue())*100.0/float64(cpuAlloc.MilliValue()))
						}
						if memAlloc.Value() > 0 {
							memUsage = fmt.Sprintf("%.2f%%", float64(memUsed.Value())*100.0/float64(memAlloc.Value()))
						}
					} else {
						if Verbose {
							fmt.Fprintf(os.Stdout, "DEBUG: No se pudieron obtener métricas para el nodo %s: %v\n", node.Name, err)
						}
					}
				}

				rows = append(rows, []string{node.Name, status, roles, kubeletVersion, cpuAllocStr, cpuUsage, memAllocStr, memUsage, age})
			}
			PrintBasicTable(headers, rows)
		}
	},
}

// init(), getNodeStatus(), getNodeRoles() se mantienen igual...
// func init() {
// 	monitorCmd.AddCommand(nodesCmd)
// }
// func getNodeStatus(node corev1.Node) string { ... }
// func getNodeRoles(node corev1.Node) string { ... }
