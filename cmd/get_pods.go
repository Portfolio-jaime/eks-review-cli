// cmd/get_pods.go
package cmd

import (
	"context"
	"encoding/json" // Para salida JSON
	"fmt"
	"os"
	"strings" // Para ToLower en outputFormat
	"time"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml" // Para salida YAML
)

// Variables para las flags de 'get pods'
var (
	podsNamespace     string
	podsAllNamespaces bool
	podsSelector      string
	podsOutputFormat  string
)

// podsGetCmd representa el comando 'monitor get pods'
var podsGetCmd = &cobra.Command{
	Use:     "pods [nombre-del-pod]",
	Aliases: []string{"po"},
	Short:   "Lista uno o más pods",
	Long: `Lista uno o más pods en el namespace actual o en todos los namespaces.
Puedes especificar un nombre de pod opcional para listar solo ese pod.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clients := GetKubeClients() // Obtiene clientes, maneja os.Exit en error

		listOptions := metav1.ListOptions{
			LabelSelector: podsSelector,
		}

		// Determinar el namespace
		// Para 'get pods', commandAllowsAllString es false (no usa "all" como valor especial para namespace)
		effectiveNamespace := GetEffectiveNamespace(podsNamespace, podsAllNamespaces, "default", false)
		if podsAllNamespaces {
			effectiveNamespace = "" // Vacío significa todos los namespaces
		}

		var podList *corev1.PodList
		var singlePod *corev1.Pod
		var err error

		if len(args) > 0 { // Si se especifica un nombre de pod
			podName := args[0]
			if effectiveNamespace == "" { // Necesario si se pide un pod específico sin -A
				// Si no se especificó un namespace y no es -A, y se busca un pod por nombre,
				// necesitamos el namespace actual o default, no ""
				// Esto es para `kubectl get pod mypod` (sin -n) vs `kubectl get pod mypod -A`
				// Si effectiveNamespace es "" por -A, está bien. Si es "" porque no se dio -n ni -A, usamos default.
				// Sin embargo, GetEffectiveNamespace ya debería haber devuelto "default" si podsNamespace era "" y podsAllNamespaces era false.
				// Así que si effectiveNamespace es "", es porque el usuario realmente quiso todos los namespaces (-A o -n all).
				// Pero `get pod <name>` no tiene sentido con `-A` sin un namespace específico para el get individual.
				// `kubectl get pod <name> -A` no es un comando válido.
				// `kubectl get pod <name> -n <ns>` sí.
				// Por ahora, si se da un nombre, y effectiveNamespace es "" (por -A), es un error o se necesita un namespace.
				// Simplificaremos: si se da un nombre de pod, se necesita un namespace efectivo (no "").
				currentContextNs := GetEffectiveNamespace(podsNamespace, false, "default", false)
				if currentContextNs == "" { // Esto no debería pasar si GetEffectiveNamespace funciona bien
					currentContextNs = "default"
				}
				if Verbose {
					fmt.Printf("DEBUG: Buscando pod '%s' en namespace '%s'\n", podName, currentContextNs)
				}
				singlePod, err = clients.Core.CoreV1().Pods(currentContextNs).Get(context.TODO(), podName, metav1.GetOptions{})
				if err != nil {
					return fmt.Errorf("error obteniendo pod '%s' en namespace '%s': %w", podName, currentContextNs, err)
				}
			} else {
				if Verbose {
					fmt.Printf("DEBUG: Buscando pod '%s' en namespace '%s'\n", podName, effectiveNamespace)
				}
				singlePod, err = clients.Core.CoreV1().Pods(effectiveNamespace).Get(context.TODO(), podName, metav1.GetOptions{})
				if err != nil {
					return fmt.Errorf("error obteniendo pod '%s' en namespace '%s': %w", podName, effectiveNamespace, err)
				}
			}
		} else { // Si no se especifica nombre de pod, listar pods
			if Verbose {
				fmt.Printf("DEBUG: Listando pods en namespace '%s' con selector '%s'\n", effectiveNamespace, podsSelector)
			}
			podList, err = clients.Core.CoreV1().Pods(effectiveNamespace).List(context.TODO(), listOptions)
			if err != nil {
				return fmt.Errorf("error listando pods: %w", err)
			}
		}

		// Procesar salida
		itemsToProcess := []corev1.Pod{}
		if singlePod != nil {
			itemsToProcess = append(itemsToProcess, *singlePod)
		} else if podList != nil {
			itemsToProcess = podList.Items
		}

		if len(itemsToProcess) == 0 {
			fmt.Fprintln(os.Stdout, "No se encontraron pods.")
			return nil
		}

		// Manejo de formatos de salida
		outputLower := strings.ToLower(podsOutputFormat)
		if outputLower == "json" {
			data, err := json.MarshalIndent(itemsToProcess, "", "  ")
			if err != nil {
				return fmt.Errorf("error convirtiendo a JSON: %w", err)
			}
			fmt.Fprintln(os.Stdout, string(data))
			return nil
		} else if outputLower == "yaml" {
			data, err := yaml.Marshal(itemsToProcess)
			if err != nil {
				return fmt.Errorf("error convirtiendo a YAML: %w", err)
			}
			fmt.Fprintln(os.Stdout, string(data))
			return nil
		}

		// Salida de tabla por defecto
		headers := []string{"NAMESPACE", "NAME", "READY", "STATUS", "RESTARTS", "AGE", "IP", "NODE"}
		if outputLower == "wide" {
			headers = append(headers, "NOMINATED NODE", "READINESS GATES")
		}

		rows := make([][]string, 0, len(itemsToProcess))
		for _, pod := range itemsToProcess {
			readyContainers := 0
			totalContainers := len(pod.Spec.Containers)
			restarts := 0
			for _, cs := range pod.Status.ContainerStatuses {
				if cs.Ready {
					readyContainers++
				}
				restarts += int(cs.RestartCount)
			}
			readyStr := fmt.Sprintf("%d/%d", readyContainers, totalContainers)
			age := metav1.Now().Sub(pod.CreationTimestamp.Time).Truncate(time.Second).String()

			row := []string{
				pod.Namespace,
				pod.Name,
				readyStr,
				string(pod.Status.Phase),
				fmt.Sprintf("%d", restarts),
				age,
				pod.Status.PodIP,
				pod.Spec.NodeName,
			}
			if outputLower == "wide" {
				row = append(row, pod.Status.NominatedNodeName, fmt.Sprintf("%v", pod.Spec.ReadinessGates)) // Simplificado
			}
			rows = append(rows, row)
		}
		PrintBasicTable(headers, rows)
		return nil
	},
}

func init() {
	getCmd.AddCommand(podsGetCmd) // Añade 'pods' como subcomando de 'monitor get'

	podsGetCmd.Flags().StringVarP(&podsNamespace, "namespace", "n", "", "Namespace para listar pods (opcional)")
	podsGetCmd.Flags().BoolVarP(&podsAllNamespaces, "all-namespaces", "A", false, "Listar pods en todos los namespaces")
	podsGetCmd.Flags().StringVarP(&podsSelector, "selector", "l", "", "Selector (label query) para filtrar pods. Ej: app=mi-app,env=prod")
	podsGetCmd.Flags().StringVarP(&podsOutputFormat, "output", "o", "", "Formato de salida. Soportado: wide, json, yaml")
}
