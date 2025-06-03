package cmd

import (
	"context"
	"fmt"
	"os"

	// "path/filepath" // Ya no es necesario aquí si GetKubeClients lo maneja
	"strings"
	"time"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/client-go/kubernetes" // Se obtiene de GetKubeClients
	// "k8s.io/client-go/tools/clientcmd" // Se usa en GetKubeClients y GetEffectiveNamespace
	// "k8s.io/client-go/util/homedir" // Se usa en GetKubeClients y GetEffectiveNamespace
)

var eventType string
var eventsNamespace string

// eventsCmd representa el comando 'monitor events'
var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Muestra eventos recientes del clúster.", // Descripción actualizada
	Long: `El comando events recupera y muestra eventos recientes de Kubernetes,
útiles para la resolución de problemas. Puedes filtrar por tipo (Warning, Normal) y namespace.`,
	Run: func(cmd *cobra.Command, args []string) {
		clients := GetKubeClients() // 1. Obtener clientes centralizadamente

		fmt.Fprintln(os.Stdout, "Recuperando eventos de Kubernetes...")

		// 2. Determinar el namespace usando la función centralizada
		// Para events, allNamespacesFlag es false, y commandAllowsAllString es true porque usa "all" como valor.
		namespaceToList := GetEffectiveNamespace(eventsNamespace, false, "default", true)

		if eventsNamespace == "" && namespaceToList == "default" && !strings.EqualFold(eventsNamespace, "all") {
			// Solo muestra este mensaje si el usuario no especificó nada y se usó "default"
			// Y si no especificó "all" explícitamente (ya que "all" también resulta en namespaceToList = "")
			fmt.Fprintf(os.Stdout, "No se especificó namespace para eventos. Usando namespace '%s'. Use -n <namespace> o -n all.\n", namespaceToList)
		} else if strings.EqualFold(eventsNamespace, "all") {
			fmt.Fprintln(os.Stdout, "Recuperando eventos de todos los namespaces.")
		} else if namespaceToList != "" {
			fmt.Fprintf(os.Stdout, "Recuperando eventos del namespace '%s'.\n", namespaceToList)
		}

		listOptions := metav1.ListOptions{}

		// Usar clients.Core
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
			var rows [][]string // Declarar rows aquí

			if Verbose { // 3. Usar la bandera Verbose
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
					fmt.Fprintln(os.Stdout, "No hay eventos que coincidan (después de filtrar).")
				}
			} else {
				rows = make([][]string, 0, len(filteredEvents)) // Inicializar rows con capacidad
				for _, event := range filteredEvents {
					lastSeen := "Desconocido"
					if !event.LastTimestamp.IsZero() {
						// 4. Usar Truncate(time.Second) para el timestamp
						duration := time.Since(event.LastTimestamp.Time).Truncate(time.Second).String()
						lastSeen = fmt.Sprintf("Hace %s", duration)
					}
					object := fmt.Sprintf("%s/%s", event.InvolvedObject.Kind, event.InvolvedObject.Name)

					rows = append(rows, []string{lastSeen, event.Type, event.Reason, object, event.Message, event.Namespace})
				}
				if Verbose { // Usar la bandera Verbose
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
