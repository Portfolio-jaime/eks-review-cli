package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	// "path/filepath" // No necesario aquí
	"strings"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/client-go/kubernetes" // Se obtiene de GetKubeClients
	// "k8s.io/client-go/tools/clientcmd" // Se usa en GetKubeClients y GetEffectiveNamespace
	// "k8s.io/client-go/util/homedir" // Se usa en GetKubeClients y GetEffectiveNamespace
)

// ... (variables de flags logPodName, logDeploymentName, etc. se mantienen) ...
// var (
// 	logPodName        string
// 	logDeploymentName string
// 	logServiceName    string
// 	logNamespace      string
// 	logContainerName  string
// 	logFollow         bool
// 	logPrevious       bool
// 	logGrep           string
// 	logTail           int64
// )

// logsCmd representa el comando 'monitor logs'
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Imprime los logs de un contenedor en un pod, deployment o servicio.",
	Long: `Imprime los logs de un contenedor en un pod, deployment o servicio.
...(resto de la descripción larga y ejemplos sin cambios)...`,
	Run: func(cmd *cobra.Command, args []string) {
		clients := GetKubeClients() // 1. Obtener clientes

		// 2. Determinar el namespace
		// Para logs, allNamespacesFlag es false, commandAllowsAllString es false (no usa "all" como valor especial)
		effectiveLogNamespace := GetEffectiveNamespace(logNamespace, false, "default", false)
		// Actualizar logNamespace para usar el valor resuelto internamente si es necesario,
		// o simplemente usar effectiveLogNamespace en las llamadas a la API.
		// Por consistencia con el resto del código, asignémoslo a logNamespace si cambió.
		if logNamespace == "" && effectiveLogNamespace == "default" {
			// Opcional: Imprimir mensaje si se usa el default y no se especificó nada
			// fmt.Fprintf(os.Stdout, "No se especificó namespace para logs. Usando namespace '%s'.\n", effectiveLogNamespace)
		}
		// Usaremos 'effectiveLogNamespace' para las llamadas a la API.

		// Validar que solo una fuente de Pod esté especificada (sin cambios)
		sourceCount := 0
		if logPodName != "" {
			sourceCount++
		}
		if logDeploymentName != "" {
			sourceCount++
		}
		if logServiceName != "" {
			sourceCount++
		}

		if sourceCount == 0 {
			fmt.Fprintf(os.Stderr, "Error: Debes especificar uno de --pod, --deployment o --service.\n")
			os.Exit(1)
		}
		if sourceCount > 1 {
			fmt.Fprintf(os.Stderr, "Error: Solo puedes especificar uno de --pod, --deployment o --service.\n")
			os.Exit(1)
		}

		var targetPodNames []string

		// Lógica para obtener Pods basada en la bandera proporcionada
		if logPodName != "" {
			targetPodNames = []string{logPodName}
		} else if logDeploymentName != "" {
			if Verbose { // 3. Usar Verbose
				fmt.Printf("DEBUG: Recuperando logs para Deployment '%s' en namespace '%s'...\n", logDeploymentName, effectiveLogNamespace)
			}
			deployment, err := clients.Core.AppsV1().Deployments(effectiveLogNamespace).Get(context.TODO(), logDeploymentName, metav1.GetOptions{})
			// ... (resto de la lógica de deployment sin cambios, usando clients.Core y effectiveLogNamespace) ...
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error obteniendo deployment '%s': %v\n", logDeploymentName, err)
				os.Exit(1)
			}
			selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error convirtiendo selector de deployment: %v\n", err)
				os.Exit(1)
			}
			pods, err := clients.Core.CoreV1().Pods(effectiveLogNamespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listando pods para deployment '%s': %v\n", logDeploymentName, err)
				os.Exit(1)
			}
			if len(pods.Items) == 0 {
				fmt.Printf("No se encontraron pods para el deployment '%s'.\n", logDeploymentName)
				os.Exit(0)
			}
			for _, p := range pods.Items {
				targetPodNames = append(targetPodNames, p.Name)
			}
		} else if logServiceName != "" {
			if Verbose { // Usar Verbose
				fmt.Printf("DEBUG: Recuperando logs para Service '%s' en namespace '%s'...\n", logServiceName, effectiveLogNamespace)
			}
			service, err := clients.Core.CoreV1().Services(effectiveLogNamespace).Get(context.TODO(), logServiceName, metav1.GetOptions{})
			// ... (resto de la lógica de service sin cambios, usando clients.Core y effectiveLogNamespace) ...
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error obteniendo service '%s': %v\n", logServiceName, err)
				os.Exit(1)
			}
			if len(service.Spec.Selector) == 0 {
				fmt.Printf("El Service '%s' no tiene selector. No se pueden encontrar pods asociados.\n", logServiceName)
				os.Exit(0)
			}
			selector := metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: service.Spec.Selector})
			pods, err := clients.Core.CoreV1().Pods(effectiveLogNamespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listando pods para service '%s': %v\n", logServiceName, err)
				os.Exit(1)
			}
			if len(pods.Items) == 0 {
				fmt.Printf("No se encontraron pods para el service '%s'.\n", logServiceName)
				os.Exit(0)
			}
			for _, p := range pods.Items {
				targetPodNames = append(targetPodNames, p.Name)
			}
		}

		logOptions := &corev1.PodLogOptions{
			Follow:    logFollow,
			Previous:  logPrevious,
			Container: logContainerName,
		}
		if logTail > 0 {
			logOptions.TailLines = &logTail
		}

		for _, podName := range targetPodNames {
			fmt.Printf("\n--- Logs para Pod: %s (Namespace: %s) ---\n", podName, effectiveLogNamespace)
			req := clients.Core.CoreV1().Pods(effectiveLogNamespace).GetLogs(podName, logOptions)
			podLogs, err := req.Stream(context.TODO())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error abriendo stream de logs para pod '%s': %v\n", podName, err)
				continue
			}
			defer podLogs.Close()

			scanner := NewLineScanner(podLogs)
			for scanner.Scan() {
				line := scanner.Text()
				if logGrep == "" || strings.Contains(line, logGrep) {
					fmt.Println(line)
				}
			}
			if err := scanner.Err(); err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error leyendo stream de logs para pod '%s': %v\n", podName, err)
			}
		}
	},
}

// init() y LineScanner se mantienen igual...
// func init() {
// 	monitorCmd.AddCommand(logsCmd)
// 	logsCmd.Flags().StringVar(&logPodName, "pod", "", "Nombre del pod del que obtener logs.")
// 	// ... resto de las flags
// }
// type LineScanner struct { ... }
// func NewLineScanner(r io.Reader) *LineScanner { ... }
// func (s *LineScanner) Scan() bool { ... }
// func (s *LineScanner) Text() string { ... }
// func (s *LineScanner) Err() error { ... }
