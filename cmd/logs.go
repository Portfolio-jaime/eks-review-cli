package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Variables para las flags del comando logs
var (
	logPodName        string
	logDeploymentName string
	logServiceName    string
	logNamespace      string
	logContainerName  string
	logFollow         bool
	logPrevious       bool
	logGrep           string // Corregido de logGrelp si era un typo
	logTail           int64
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Imprime los logs de un contenedor en un pod, deployment o servicio.",
	Long: `Imprime los logs de un contenedor en un pod, deployment o servicio.
...(ejemplos)...`,
	Run: func(cmd *cobra.Command, args []string) {
		clients, err := GetKubeClients()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creando clientes de Kubernetes: %v\n", err)
			os.Exit(1)
		}

		effectiveLogNamespace := GetEffectiveNamespace(logNamespace, false, "default", false)

		if logNamespace == "" && effectiveLogNamespace == "default" {
			// Opcional: fmt.Fprintf(os.Stdout, "No se especificó namespace para logs. Usando namespace '%s'.\n", effectiveLogNamespace)
		}

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
		if logPodName != "" {
			targetPodNames = []string{logPodName}
		} else if logDeploymentName != "" {
			if Verbose {
				fmt.Printf("DEBUG: Recuperando logs para Deployment '%s' en namespace '%s'...\n", logDeploymentName, effectiveLogNamespace)
			}
			deployment, err := clients.Core.AppsV1().Deployments(effectiveLogNamespace).Get(context.TODO(), logDeploymentName, metav1.GetOptions{})
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
			if Verbose {
				fmt.Printf("DEBUG: Recuperando logs para Service '%s' en namespace '%s'...\n", logServiceName, effectiveLogNamespace)
			}
			service, err := clients.Core.CoreV1().Services(effectiveLogNamespace).Get(context.TODO(), logServiceName, metav1.GetOptions{})
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
			Follow: logFollow, Previous: logPrevious, Container: logContainerName,
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

func init() {
	monitorCmd.AddCommand(logsCmd)
	logsCmd.Flags().StringVar(&logPodName, "pod", "", "Nombre del pod del que obtener logs.")
	logsCmd.Flags().StringVar(&logDeploymentName, "deployment", "", "Nombre del deployment del que obtener logs.")
	logsCmd.Flags().StringVar(&logServiceName, "service", "", "Nombre del service del que obtener logs.")
	logsCmd.Flags().StringVarP(&logNamespace, "namespace", "n", "", "Si está presente, el ámbito del namespace para esta solicitud CLI.")
	logsCmd.Flags().StringVarP(&logContainerName, "container", "c", "", "Nombre del contenedor.")
	logsCmd.Flags().BoolVarP(&logFollow, "follow", "f", false, "Especificar si los logs deben ser transmitidos.")
	logsCmd.Flags().BoolVarP(&logPrevious, "previous", "p", false, "Si es true, imprime los logs de la instancia previa del contenedor.")
	logsCmd.Flags().StringVar(&logGrep, "grep", "", "Filtrar logs por una cadena de texto.")
	logsCmd.Flags().Int64Var(&logTail, "tail", -1, "Líneas desde el final de los logs a mostrar.")
}

// LineScanner y sus métodos (deben estar aquí)
type LineScanner struct { /* ... */
	reader io.Reader
	buf    []byte
	err    error
	eof    bool
}

func NewLineScanner(r io.Reader) *LineScanner {
	return &LineScanner{reader: r, buf: make([]byte, 0, 4*1024)}
}
func (s *LineScanner) Scan() bool {
	if s.err != nil || s.eof {
		return false
	}
	for {
		n, readErr := s.reader.Read(s.buf[len(s.buf):cap(s.buf)])
		if n > 0 {
			s.buf = s.buf[:len(s.buf)+n]
		}
		if lineEnd := strings.IndexByte(string(s.buf), '\n'); lineEnd >= 0 {
			s.err = nil
			return true
		}
		if readErr == io.EOF {
			s.eof = true
			if len(s.buf) > 0 {
				s.err = nil
				return true
			}
			s.err = io.EOF
			return false
		}
		if readErr != nil {
			s.err = readErr
			return false
		}
		if len(s.buf) == cap(s.buf) {
			newBuf := make([]byte, len(s.buf), len(s.buf)*2)
			copy(newBuf, s.buf)
			s.buf = newBuf
		}
	}
}
func (s *LineScanner) Text() string {
	if lineEnd := strings.IndexByte(string(s.buf), '\n'); lineEnd >= 0 {
		line := s.buf[:lineEnd]
		s.buf = s.buf[lineEnd+1:]
		return string(line)
	}
	line := s.buf
	s.buf = s.buf[:0]
	return string(line)
}
func (s *LineScanner) Err() error { return s.err }
