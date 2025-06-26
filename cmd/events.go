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

var eventType string
var eventsNamespace string

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Muestra eventos recientes del clúster.",
	Long: `El comando events recupera y muestra eventos recientes de Kubernetes,
útiles para la resolución de problemas. Puedes filtrar por tipo (Warning, Normal) y namespace.`,
	Run: func(cmd *cobra.Command, args []string) {
		clients, err := GetKubeClients()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creando clientes de Kubernetes: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stdout, "Recuperando eventos de Kubernetes...")

		namespaceToList := GetEffectiveNamespace(eventsNamespace, false, "default", true)

		if eventsNamespace == "" && namespaceToList == "default" && !strings.EqualFold(eventsNamespace, "all") {
			fmt.Fprintf(os.Stdout, "No se especificó namespace para eventos. Usando namespace '%s'. Use -n <namespace> o -n all.\n", namespaceToList)
		} else if strings.EqualFold(eventsNamespace, "all") {
			fmt.Fprintln(os.Stdout, "Recuperando eventos de todos los namespaces.")
		} else if namespaceToList != "" { // Añadido para ser explícito
			fmt.Fprintf(os.Stdout, "Recuperando eventos del namespace '%s'.\n", namespaceToList)
		}

		listOptions := metav1.ListOptions{}
		events, err := clients.Core.CoreV1().Events(namespaceToList).List(context.TODO(), listOptions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listando eventos: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stdout, "\n--- Eventos ---")
		if len(events.Items) == 0 {
			fmt.Fprintln(os.Stdout, "No se encontraron eventos.")
		} else {
			headers := []string{"ÚLTIMA VEZ", "TIPO", "RAZÓN", "OBJETO", "MENSAJE", "NAMESPACE"}
			var rows [][]string

			if Verbose {
				fmt.Fprintf(os.Stdout, "DEBUG: Eventos encontrados (crudos): %d\n", len(events.Items))
			}

			filteredEvents := []corev1.Event{}
			for _, event := range events.Items {
				if eventType == "" || strings.EqualFold(event.Type, eventType) {
					filteredEvents = append(filteredEvents, event)
				}
			}

			if len(filteredEvents) == 0 {
				if eventType != "" {
					fmt.Fprintf(os.Stdout, "No se encontraron eventos de tipo '%s'.\n", eventType)
				} else {
					// Este mensaje podría ser confuso si events.Items no estaba vacío pero filteredEvents sí.
					// Mejor ser más específico o simplemente no imprimir nada si no hay filtro de tipo.
					fmt.Fprintln(os.Stdout, "No hay eventos que coincidan (después de filtrar).")
				}
			} else {
				rows = make([][]string, 0, len(filteredEvents))
				for _, event := range filteredEvents {
					lastSeen := "Desconocido"
					if !event.LastTimestamp.IsZero() {
						duration := time.Since(event.LastTimestamp.Time).Truncate(time.Second).String()
						lastSeen = fmt.Sprintf("Hace %s", duration)
					}
					object := fmt.Sprintf("%s/%s", event.InvolvedObject.Kind, event.InvolvedObject.Name)
					rows = append(rows, []string{lastSeen, event.Type, event.Reason, object, event.Message, event.Namespace})
				}
				if Verbose {
					fmt.Fprintf(os.Stdout, "DEBUG: Eventos filtrados añadidos a la tabla. Recuento final de filas: %d\n", len(rows))
				}
				PrintBasicTable(headers, rows)
			}
		}
	},
}

func init() {
	monitorCmd.AddCommand(eventsCmd)
	eventsCmd.Flags().StringVarP(&eventType, "type", "T", "", "Filtrar eventos por tipo (ej., 'Warning', 'Normal'). No sensible a mayúsculas.")
	eventsCmd.Flags().StringVarP(&eventsNamespace, "namespace", "n", "", "Si está presente, el ámbito del namespace para esta solicitud CLI. Usa 'all' para todos los namespaces.")
}
