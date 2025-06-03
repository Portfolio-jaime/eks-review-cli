package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// Asegúrate de tener estas importaciones también para las nuevas funciones:
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	// "k8s.io/client-go/rest" // Descomenta si decides devolver Config en KubeClients
)

// PrintBasicTable imprime una tabla formateada de forma básica en la consola.
// Adapta automáticamente el ancho de las columnas.
func PrintBasicTable(headers []string, rows [][]string) {
	// ... (tu código existente de PrintBasicTable se mantiene igual)
	// ... (me aseguro de que esto no se pierda)
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
			formatString += "  " // Espacio de 2 caracteres entre columnas
			separator += "--"
		}
	}
	formatString += "\n"
	separator += "\n"

	// Imprimir encabezados
	headerArgs := make([]interface{}, len(headers))
	for i, h := range headers {
		headerArgs[i] = h
	}
	fmt.Printf(formatString, headerArgs...)
	fmt.Printf(separator)

	// Imprimir filas
	for _, row := range rows {
		rowArgs := make([]interface{}, len(row))
		for i, cell := range row {
			rowArgs[i] = cell
		}
		fmt.Printf(formatString, rowArgs...)
	}
	fmt.Println() // Línea vacía al final de la tabla
}

// --- NUEVAS FUNCIONES CENTRALIZADAS ---

// KubeClients contiene los clientes de Kubernetes inicializados.
type KubeClients struct {
	Core    kubernetes.Interface
	Metrics metrics.Interface
	// Config *rest.Config // Opcionalmente, devuelve la configuración cruda
}

// GetKubeClients inicializa y devuelve los clientes core y de métricas de Kubernetes.
// Sale del programa si la creación del cliente core falla.
// La creación del cliente de métricas es opcional; se imprime una advertencia si falla.
func GetKubeClients() *KubeClients {
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error construyendo kubeconfig: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creando cliente de Kubernetes: %v\n", err)
		os.Exit(1)
	}

	var metricsClientset metrics.Interface
	metricsClientset, err = metrics.NewForConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Advertencia: No se pudo crear el cliente de métricas (el servidor de métricas podría no estar instalado): %v\n", err)
		metricsClientset = nil
	}

	return &KubeClients{
		Core:    clientset,
		Metrics: metricsClientset,
	}
}

// GetEffectiveNamespace determina el namespace a utilizar.
//
// Parámetros:
//   - namespaceFlag: El valor de la bandera --namespace para el comando.
//   - allNamespacesFlag: Un booleano que indica si se usó una bandera estilo --all-namespaces.
//   - defaultNamespaceVal: El namespace al que se recurrirá si no se encuentra ninguno.
//   - commandAllowsAllString: true si el comando usa "all" como valor de cadena para su bandera de namespace.
//
// Devuelve la cadena del namespace. Una cadena vacía típicamente significa "todos los namespaces".
func GetEffectiveNamespace(namespaceFlag string, allNamespacesFlag bool, defaultNamespaceVal string, commandAllowsAllString bool) string {
	if allNamespacesFlag {
		return "" // "" significa todos los namespaces para client-go
	}

	if namespaceFlag != "" {
		if commandAllowsAllString && strings.ToLower(namespaceFlag) == "all" {
			return ""
		}
		return namespaceFlag
	}

	// Intenta obtener del contexto actual si no se proporcionó una bandera de namespace específica
	// No es necesario pasar kubeconfigPath aquí si siempre usamos el default para esta lógica
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	// Carga no interactiva para obtener el namespace del contexto actual
	// Advertencia: Esta parte puede ser un poco lenta si se llama muchas veces y el kubeconfig es complejo.
	// Para CLIs, usualmente es aceptable.
	clientConfigLoadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	clientConfigLoadingRules.ExplicitPath = kubeconfigPath
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientConfigLoadingRules, configOverrides)

	currentNamespace, _, err := kubeConfig.Namespace()
	if err == nil && currentNamespace != "" {
		return currentNamespace
	}

	// Fallback al predeterminado
	if defaultNamespaceVal == "" {
		defaultNamespaceVal = "default"
	}
	return defaultNamespaceVal
}
