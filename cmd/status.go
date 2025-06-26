package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1" // Asegúrate que esta importación esté
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes" // Necesario para el tipo del parámetro clientset
)

var allNamespaces bool
var targetNamespace string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Proporciona un resumen rápido del estado de los recursos.",
	Long: `El comando status recupera un resumen de Pods, Deployments, Services,
e Ingresses en un namespace dado o en todos los namespaces.`,
	Run: func(cmd *cobra.Command, args []string) {
		clients, err := GetKubeClients()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creando clientes de Kubernetes: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Recuperando estado de recursos de Kubernetes...")
		namespaceToList := GetEffectiveNamespace(targetNamespace, allNamespaces, "default", false)

		if !allNamespaces && targetNamespace == "" && namespaceToList == "default" {
			fmt.Fprintf(os.Stdout, "No se especificó namespace. Usando namespace '%s'. Use -n <namespace> o -A / --all-namespaces.\n", namespaceToList)
		} else if allNamespaces {
			fmt.Fprintln(os.Stdout, "Recuperando estado de recursos de todos los namespaces.")
		} else if namespaceToList != "" {
			fmt.Fprintf(os.Stdout, "Recuperando estado de recursos del namespace '%s'.\n", namespaceToList)
		}

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

func listPods(clientset kubernetes.Interface, namespace string) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listando pods: %v\n", err)
		return
	}

	fmt.Fprintln(os.Stdout, "\n--- Pods ---")
	if len(pods.Items) == 0 {
		fmt.Fprintln(os.Stdout, "No se encontraron pods.")
		return
	}

	headers := []string{"NOMBRE", "NAMESPACE", "ESTADO", "REINICIOS", "EDAD"}
	rows := make([][]string, 0, len(pods.Items))
	if Verbose {
		fmt.Fprintf(os.Stdout, "DEBUG: Pods encontrados: %d\n", len(pods.Items))
	}

	for _, pod := range pods.Items {
		restarts := 0
		for _, cs := range pod.Status.ContainerStatuses {
			restarts += int(cs.RestartCount)
		}
		age := metav1.Now().Sub(pod.CreationTimestamp.Time).Truncate(time.Second).String()
		rows = append(rows, []string{pod.Name, pod.Namespace, string(pod.Status.Phase), fmt.Sprintf("%d", restarts), age})
	}
	PrintBasicTable(headers, rows)
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
		return
	}

	headers := []string{"NOMBRE", "NAMESPACE", "LISTOS", "ACTUALIZADOS", "DISPONIBLES", "EDAD"}
	rows := make([][]string, 0, len(deployments.Items))
	if Verbose {
		fmt.Fprintf(os.Stdout, "DEBUG: Deployments encontrados: %d\n", len(deployments.Items))
	}

	for _, deploy := range deployments.Items {
		readyReplicas := int32(0)
		if deploy.Spec.Replicas != nil { // deploy.Spec.Replicas es un puntero
			readyReplicas = *deploy.Spec.Replicas
		}
		ready := fmt.Sprintf("%d/%d", deploy.Status.ReadyReplicas, readyReplicas)
		upToDate := fmt.Sprintf("%d", deploy.Status.UpdatedReplicas)
		available := fmt.Sprintf("%d", deploy.Status.AvailableReplicas)
		age := metav1.Now().Sub(deploy.CreationTimestamp.Time).Truncate(time.Second).String()
		rows = append(rows, []string{deploy.Name, deploy.Namespace, ready, upToDate, available, age})
	}
	PrintBasicTable(headers, rows)
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
		return
	}

	headers := []string{"NOMBRE", "NAMESPACE", "TIPO", "CLUSTER-IP", "IP-EXTERNA", "PUERTO(S)", "EDAD"}
	rows := make([][]string, 0, len(services.Items))
	if Verbose {
		fmt.Fprintf(os.Stdout, "DEBUG: Services encontrados: %d\n", len(services.Items))
	}

	for _, svc := range services.Items {
		externalIP := "<none>"
		if svc.Spec.Type == corev1.ServiceTypeLoadBalancer { // Usa corev1 aquí
			if len(svc.Status.LoadBalancer.Ingress) > 0 {
				if svc.Status.LoadBalancer.Ingress[0].IP != "" {
					externalIP = svc.Status.LoadBalancer.Ingress[0].IP
				} else if svc.Status.LoadBalancer.Ingress[0].Hostname != "" {
					externalIP = svc.Status.LoadBalancer.Ingress[0].Hostname
				} else {
					externalIP = "<pending>"
				}
			} else {
				externalIP = "<pending>"
			}
		} else if len(svc.Spec.ExternalIPs) > 0 {
			externalIP = strings.Join(svc.Spec.ExternalIPs, ",")
		}

		var portStrings []string
		for _, port := range svc.Spec.Ports {
			pStr := fmt.Sprintf("%d", port.Port)
			if port.NodePort > 0 {
				pStr += fmt.Sprintf(":%d", port.NodePort)
			}
			pStr += fmt.Sprintf("/%s", port.Protocol)
			portStrings = append(portStrings, pStr)
		}
		age := metav1.Now().Sub(svc.CreationTimestamp.Time).Truncate(time.Second).String()
		rows = append(rows, []string{svc.Name, svc.Namespace, string(svc.Spec.Type), svc.Spec.ClusterIP, externalIP, strings.Join(portStrings, ","), age})
	}
	PrintBasicTable(headers, rows)
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
		return
	}

	headers := []string{"NOMBRE", "NAMESPACE", "CLASE", "HOSTS", "DIRECCIÓN", "PUERTOS", "EDAD"}
	rows := make([][]string, 0, len(ingresses.Items))
	if Verbose {
		fmt.Fprintf(os.Stdout, "DEBUG: Ingresses encontrados: %d\n", len(ingresses.Items))
	}

	for _, ingress := range ingresses.Items {
		address := "<none>"
		if len(ingress.Status.LoadBalancer.Ingress) > 0 {
			var addresses []string
			for _, ingStatus := range ingress.Status.LoadBalancer.Ingress {
				if ingStatus.IP != "" {
					addresses = append(addresses, ingStatus.IP)
				}
				if ingStatus.Hostname != "" {
					addresses = append(addresses, ingStatus.Hostname)
				}
			}
			if len(addresses) > 0 {
				address = strings.Join(addresses, ",")
			} else {
				address = "<pending>"
			}
		}

		var hosts []string
		for _, rule := range ingress.Spec.Rules {
			hosts = append(hosts, rule.Host)
		}

		className := ""
		if ingress.Spec.IngressClassName != nil {
			className = *ingress.Spec.IngressClassName
		}

		age := metav1.Now().Sub(ingress.CreationTimestamp.Time).Truncate(time.Second).String()

		// Para PUERTOS en Ingress, es comúnmente 80/443, manejado por el controller.
		// Es difícil extraer esto de forma genérica del objeto Ingress.
		portStr := "80, 443" // Placeholder o puedes intentar lógica más compleja.

		rows = append(rows, []string{ingress.Name, ingress.Namespace, className, strings.Join(hosts, ","), address, portStr, age})
	}
	PrintBasicTable(headers, rows)
}
