package cmd

import (
	"context"
	"fmt"
	"os"

	// "path/filepath" // No necesario aquí
	"strings"
	"time" // Necesario para age

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/kubernetes" // Se obtiene de GetKubeClients
	// "k8s.io/client-go/tools/clientcmd" // Se usa en GetKubeClients y GetEffectiveNamespace
	// "k8s.io/client-go/util/homedir" // Se usa en GetKubeClients y GetEffectiveNamespace
)

var allNamespaces bool
var targetNamespace string

// statusCmd representa el comando 'monitor status'
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Proporciona un resumen rápido del estado de los recursos.",
	Long: `El comando status recupera un resumen de Pods, Deployments, Services,
Ingresses en un namespace dado o en todos los namespaces.`, // Eliminé PVs ya que no se listan actualmente
	Run: func(cmd *cobra.Command, args []string) {
		clients := GetKubeClients() // 1. Obtener clientes

		fmt.Println("Recuperando estado de recursos de Kubernetes...")

		// 2. Determinar el namespace
		// Para status, allNamespacesFlag es la variable 'allNamespaces', commandAllowsAllString es false.
		namespaceToList := GetEffectiveNamespace(targetNamespace, allNamespaces, "default", false)

		if !allNamespaces && targetNamespace == "" && namespaceToList == "default" {
			fmt.Fprintf(os.Stdout, "No se especificó namespace. Usando namespace '%s'. Use -n <namespace> o -A / --all-namespaces.\n", namespaceToList)
		} else if allNamespaces {
			fmt.Fprintln(os.Stdout, "Recuperando estado de recursos de todos los namespaces.")
		} else if namespaceToList != "" {
			fmt.Fprintf(os.Stdout, "Recuperando estado de recursos del namespace '%s'.\n", namespaceToList)
		}

		// Pasar clients.Core a las funciones de listado
		listPods(clients.Core, namespaceToList)
		listDeployments(clients.Core, namespaceToList)
		listServices(clients.Core, namespaceToList)
		listIngresses(clients.Core, namespaceToList)
	},
}

func init() {
	monitorCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "Si es true, lista el/los objeto(s) solicitado(s) en todos los namespaces.")
	statusCmd.Flags().StringVarP(&targetNamespace, "namespace", "n", "", "Si está presente, el ámbito del namespace para esta solicitud CLI.")
}

// Modificar las funciones listX para aceptar kubernetes.Interface y usar Verbose
func listPods(clientset kubernetes.Interface, namespace string) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listando pods: %v\n", err)
		return
	}

	fmt.Fprintln(os.Stdout, "\n--- Pods ---")
	if len(pods.Items) == 0 {
		fmt.Fprintln(os.Stdout, "No se encontraron pods.")
	} else {
		headers := []string{"NOMBRE", "NAMESPACE", "ESTADO", "REINICIOS", "EDAD"}
		rows := make([][]string, 0, len(pods.Items))

		if Verbose { // 3. Usar Verbose
			fmt.Fprintf(os.Stdout, "DEBUG: Pods encontrados: %d\n", len(pods.Items))
		}

		for _, pod := range pods.Items {
			restarts := 0
			for _, containerStatus := range pod.Status.ContainerStatuses {
				restarts += int(containerStatus.RestartCount)
			}
			age := metav1.Now().Sub(pod.CreationTimestamp.Time).Truncate(time.Second).String() // Usar Truncate

			rows = append(rows, []string{pod.Name, pod.Namespace, string(pod.Status.Phase), fmt.Sprintf("%d", restarts), age})
		}
		PrintBasicTable(headers, rows)
	}
}

func listDeployments(clientset kubernetes.Interface, namespace string) {
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listando deployments: %v\n", err)
		return
	}

	fmt.Fprintln(os.Stdout, "\n--- Deployments ---")
	if len(deployments.Items) == 0 {
		fmt.Fprintln(os.Stdout, "No se encontraron deployments.")
	} else {
		headers := []string{"NOMBRE", "NAMESPACE", "LISTOS", "ACTUALIZADOS", "DISPONIBLES", "EDAD"}
		rows := make([][]string, 0, len(deployments.Items))

		if Verbose { // Usar Verbose
			fmt.Fprintf(os.Stdout, "DEBUG: Deployments encontrados: %d\n", len(deployments.Items))
		}

		for _, deploy := range deployments.Items {
			ready := fmt.Sprintf("%d/%d", deploy.Status.ReadyReplicas, deploy.Spec.Replicas) // Usar ReadyReplicas y Spec.Replicas
			upToDate := fmt.Sprintf("%d", deploy.Status.UpdatedReplicas)
			available := fmt.Sprintf("%d", deploy.Status.AvailableReplicas)
			age := metav1.Now().Sub(deploy.CreationTimestamp.Time).Truncate(time.Second).String() // Usar Truncate

			rows = append(rows, []string{deploy.Name, deploy.Namespace, ready, upToDate, available, age})
		}
		PrintBasicTable(headers, rows)
	}
}

func listServices(clientset kubernetes.Interface, namespace string) {
	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listando services: %v\n", err)
		return
	}

	fmt.Fprintln(os.Stdout, "\n--- Services ---")
	if len(services.Items) == 0 {
		fmt.Fprintln(os.Stdout, "No se encontraron services.")
	} else {
		headers := []string{"NOMBRE", "NAMESPACE", "TIPO", "CLUSTER-IP", "IP-EXTERNA", "PUERTO(S)", "EDAD"}
		rows := make([][]string, 0, len(services.Items))

		if Verbose { // Usar Verbose
			fmt.Fprintf(os.Stdout, "DEBUG: Services encontrados: %d\n", len(services.Items))
		}

		for _, svc := range services.Items {
			externalIP := "<none>" // Default a <none> si no hay IP externa
			if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
				if len(svc.Status.LoadBalancer.Ingress) > 0 {
					externalIP = svc.Status.LoadBalancer.Ingress[0].IP
					if externalIP == "" { // A veces es Hostname en lugar de IP
						externalIP = svc.Status.LoadBalancer.Ingress[0].Hostname
					}
				} else {
					externalIP = "<pending>" // Si es LoadBalancer pero aún no tiene IP
				}
			} else if len(svc.Spec.ExternalIPs) > 0 {
				externalIP = strings.Join(svc.Spec.ExternalIPs, ",")
			}

			ports := []string{}
			for _, port := range svc.Spec.Ports {
				ports = append(ports, fmt.Sprintf("%d:%d/%s", port.Port, port.NodePort, port.Protocol)) // Incluir NodePort si existe
			}
			age := metav1.Now().Sub(svc.CreationTimestamp.Time).Truncate(time.Second).String() // Usar Truncate

			rows = append(rows, []string{svc.Name, svc.Namespace, string(svc.Spec.Type), svc.Spec.ClusterIP, externalIP, strings.Join(ports, ", "), age})
		}
		PrintBasicTable(headers, rows)
	}
}

func listIngresses(clientset kubernetes.Interface, namespace string) {
	ingresses, err := clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listando ingresses: %v\n", err)
		return
	}

	fmt.Fprintln(os.Stdout, "\n--- Ingresses ---")
	if len(ingresses.Items) == 0 {
		fmt.Fprintln(os.Stdout, "No se encontraron ingresses.")
	} else {
		headers := []string{"NOMBRE", "NAMESPACE", "CLASE", "HOSTS", "DIRECCIÓN", "PUERTOS", "EDAD"} // Agregué PUERTOS
		rows := make([][]string, 0, len(ingresses.Items))

		if Verbose { // Usar Verbose
			fmt.Fprintf(os.Stdout, "DEBUG: Ingresses encontrados: %d\n", len(ingresses.Items))
		}

		for _, ingress := range ingresses.Items {
			address := "<none>"
			ports := "" // Para Ingress, los puertos suelen ser 80, 443 implícitos o definidos en el Ingress Controller.
			// Esta información es más compleja de obtener de forma genérica desde el objeto Ingress en sí.
			// Podrías mostrar "80, 443" o dejarlo como N/A.
			// Por ahora, lo dejaré simple.

			if len(ingress.Status.LoadBalancer.Ingress) > 0 {
				ips := []string{}
				for _, ingStatus := range ingress.Status.LoadBalancer.Ingress {
					if ingStatus.IP != "" {
						ips = append(ips, ingStatus.IP)
					} else if ingStatus.Hostname != "" {
						ips = append(ips, ingStatus.Hostname)
					}
				}
				if len(ips) > 0 {
					address = strings.Join(ips, ",")
				} else {
					address = "<pending>"
				}
			}

			hosts := []string{}
			for _, rule := range ingress.Spec.Rules {
				hosts = append(hosts, rule.Host)
			}
			className := ""
			if ingress.Spec.IngressClassName != nil {
				className = *ingress.Spec.IngressClassName
			}
			age := metav1.Now().Sub(ingress.CreationTimestamp.Time).Truncate(time.Second).String() // Usar Truncate

			// Para los puertos del Ingress: típicamente 80/443, gestionado por el Ingress Controller
			// Aquí asumiremos los más comunes o lo dejaremos vacío.
			// Podrías mostrar "http/https" o algo similar si lo deseas.
			// Para simplificar, usaré "80, 443" como placeholder.
			// O puedes intentar obtenerlo de las reglas si definen backends con puertos, aunque es más complejo.
			portStr := "80, 443" // Placeholder

			rows = append(rows, []string{ingress.Name, ingress.Namespace, className, strings.Join(hosts, ", "), address, portStr, age})
		}
		PrintBasicTable(headers, rows)
	}
}
