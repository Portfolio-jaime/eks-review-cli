package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	// corev1 "k8s.io/api/core/v1" // Importa si alguna función aquí lo necesita
)

// KubeClients contiene los clientes de Kubernetes inicializados.
type KubeClients struct {
	Core    kubernetes.Interface
	Metrics metrics.Interface
}

// GetKubeClients inicializa y devuelve los clientes core y de métricas de Kubernetes.
// Devuelve un error en lugar de finalizar el programa para permitir un mejor manejo
// de fallos y facilitar las pruebas unitarias de los comandos.
func GetKubeClients() (*KubeClients, error) {
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("construyendo kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creando cliente de Kubernetes: %w", err)
	}

	var metricsClientset metrics.Interface
	metricsErr := (error)(nil) // Para claridad
	metricsClientset, metricsErr = metrics.NewForConfig(config)
	if metricsErr != nil {
		fmt.Fprintf(os.Stderr, "Advertencia: No se pudo crear el cliente de métricas (el servidor de métricas podría no estar instalado): %v\n", metricsErr)
		metricsClientset = nil
	}

	return &KubeClients{
		Core:    clientset,
		Metrics: metricsClientset,
	}, nil
}

// GetEffectiveNamespace determina el namespace a utilizar.
func GetEffectiveNamespace(namespaceFlag string, allNamespacesFlag bool, defaultNamespaceVal string, commandAllowsAllString bool) string {
	if allNamespacesFlag {
		return ""
	}

	if namespaceFlag != "" {
		if commandAllowsAllString && strings.ToLower(namespaceFlag) == "all" {
			return ""
		}
		return namespaceFlag
	}

	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	clientConfigLoadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	clientConfigLoadingRules.ExplicitPath = kubeconfigPath
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientConfigLoadingRules, configOverrides)

	currentNamespace, _, err := kubeConfig.Namespace()
	if err == nil && currentNamespace != "" {
		return currentNamespace
	}

	if defaultNamespaceVal == "" {
		defaultNamespaceVal = "default"
	}
	return defaultNamespaceVal
}

// PrintBasicTable imprime una tabla formateada de forma básica en la consola.
func PrintBasicTable(headers []string, rows [][]string) {
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}

	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	formatString := ""
	separator := ""
	for i, width := range colWidths {
		formatString += fmt.Sprintf("%%-%ds", width)
		separator += strings.Repeat("-", width)
		if i < len(colWidths)-1 {
			formatString += "  "
			separator += "--"
		}
	}
	formatString += "\n"
	separator += "\n"

	headerArgs := make([]interface{}, len(headers))
	for i, h := range headers {
		headerArgs[i] = h
	}
	fmt.Print(fmt.Sprintf(formatString, headerArgs...))
	fmt.Print(separator)

	for _, row := range rows {
		rowArgs := make([]interface{}, len(row))
		for i, cell := range row {
			rowArgs[i] = cell
		}
		fmt.Print(fmt.Sprintf(formatString, rowArgs...))
	}
	fmt.Println()
}
