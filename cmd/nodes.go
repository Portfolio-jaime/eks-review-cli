package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Muestra información detallada sobre los nodos del clúster.",
	Long: `El comando nodes recupera y muestra información detallada sobre
los nodos del clúster de Kubernetes, incluyendo su estado, roles y uso de recursos.
El uso de recursos requiere un servidor de métricas instalado en el clúster.`,
	Run: func(cmd *cobra.Command, args []string) {
		clients, err := GetKubeClients()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creando clientes de Kubernetes: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stdout, "Recuperando información de nodos de Kubernetes...")

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
			rows := make([][]string, 0, len(nodes.Items))

			if Verbose {
				fmt.Fprintf(os.Stdout, "DEBUG: Nodos encontrados: %d\n", len(nodes.Items))
			}

			for _, node := range nodes.Items {
				status := getNodeStatus(node) // Asegúrate que esta función exista
				roles := getNodeRoles(node)   // Asegúrate que esta función exista
				kubeletVersion := node.Status.NodeInfo.KubeletVersion
				age := metav1.Now().Sub(node.CreationTimestamp.Time).Truncate(time.Second).String()

				cpuAlloc := node.Status.Allocatable[corev1.ResourceCPU]
				memAlloc := node.Status.Allocatable[corev1.ResourceMemory]
				cpuAllocStr := cpuAlloc.String()
				memAllocStr := memAlloc.String()

				cpuUsage := "N/A"
				memUsage := "N/A"

				if clients.Metrics != nil {
					nodeMetrics, errMetrics := clients.Metrics.MetricsV1beta1().NodeMetricses().Get(context.TODO(), node.Name, metav1.GetOptions{})
					if errMetrics == nil {
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
							fmt.Fprintf(os.Stdout, "DEBUG: No se pudieron obtener métricas para el nodo %s: %v\n", node.Name, errMetrics)
						}
					}
				}
				rows = append(rows, []string{node.Name, status, roles, kubeletVersion, cpuAllocStr, cpuUsage, memAllocStr, memUsage, age})
			}
			PrintBasicTable(headers, rows)
		}
	},
}

func init() {
	monitorCmd.AddCommand(nodesCmd)
}

// Helper functions (asegúrate que estén aquí)
func getNodeStatus(node corev1.Node) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			if condition.Status == corev1.ConditionTrue {
				return "Ready"
			}
			return "NotReady"
		}
	}
	return "Unknown"
}

func getNodeRoles(node corev1.Node) string {
	var roles []string
	for k, v := range node.Labels {
		if strings.HasPrefix(k, "node-role.kubernetes.io/") {
			role := strings.TrimPrefix(k, "node-role.kubernetes.io/")
			if role != "" { // Añadir solo si el rol no es vacío
				roles = append(roles, role)
			}
		}
		// Considerar también la etiqueta "kubernetes.io/role" si es relevante para tu entorno
		if k == "kubernetes.io/role" && v != "" {
			// Evitar duplicados si ya se añadió por el prefijo
			isDup := false
			for _, r := range roles {
				if r == v {
					isDup = true
					break
				}
			}
			if !isDup {
				roles = append(roles, v)
			}
		}
	}
	if len(roles) == 0 {
		return "<none>"
	}
	return strings.Join(roles, ",")
}
