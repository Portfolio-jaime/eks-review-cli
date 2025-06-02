package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings" // Importa para strings.Join

	"github.com/olekukonko/tablewriter" // Importa la librería para tablas
	"github.com/spf13/cobra"            // Para Deployments
	// Para Pods y Services
	// Para Ingresses
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var allNamespaces bool
var targetNamespace string

// statusCmd representa el comando 'monitor status'
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Provides a quick summary of resource states.",
	Long: `The status command retrieves a summary of Pods, Deployments, Services,
Ingresses, and Persistent Volumes in a given namespace or across all namespaces.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Construir la ruta al kubeconfig
		kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

		// Cargar la configuración de Kubernetes
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error building kubeconfig: %v\n", err)
			os.Exit(1)
		}

		// Crear un cliente de Kubernetes
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating Kubernetes client: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Fetching Kubernetes resource status...")

		// Determinar el namespace objetivo
		var namespaceToList string
		if allNamespaces {
			namespaceToList = "" // Una cadena vacía significa todos los namespaces
		} else if targetNamespace != "" {
			namespaceToList = targetNamespace
		} else {
			// Si no se especifica --all-namespaces ni -n, intenta obtener el namespace del contexto actual
			rawConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
				&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
				&clientcmd.ConfigOverrides{},
			).RawConfig()
			if err == nil && rawConfig.CurrentContext != "" {
				currentContext := rawConfig.Contexts[rawConfig.CurrentContext]
				if currentContext.Namespace != "" {
					namespaceToList = currentContext.Namespace
				}
			}
			if namespaceToList == "" {
				namespaceToList = "default" // Por defecto a "default" si no se encuentra un namespace en el contexto
				fmt.Printf("No namespace specified. Using '%s' namespace. Use -n <namespace> or --all-namespaces.\n", namespaceToList)
			}
		}

		// --- Listado de Pods ---
		listPods(clientset, namespaceToList)

		// --- Listado de Deployments ---
		listDeployments(clientset, namespaceToList)

		// --- Listado de Services ---
		listServices(clientset, namespaceToList)

		// --- Listado de Ingresses ---
		listIngresses(clientset, namespaceToList)

		// TODO: Puedes añadir Persistent Volumes, StatefulSets, DaemonSets, etc. siguiendo el mismo patrón.
	},
}

// init function for the status command
func init() {
	monitorCmd.AddCommand(statusCmd)

	// Define las banderas (flags) para el comando status
	statusCmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If true, list the requested object(s) across all namespaces.")
	statusCmd.Flags().StringVarP(&targetNamespace, "namespace", "n", "", "If present, the namespace scope for this CLI request.")
}

// listPods function to get and display pod status
func listPods(clientset *kubernetes.Clientset, namespace string) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing pods: %v\n", err)
		// No salir, para que otros listados puedan continuar
		return
	}

	fmt.Println("\n--- Pods ---")
	if len(pods.Items) == 0 {
		fmt.Println("No pods found.")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"NAME", "NAMESPACE", "STATUS", "RESTARTS", "AGE"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")

		for _, pod := range pods.Items {
			restarts := 0
			for _, containerStatus := range pod.Status.ContainerStatuses {
				restarts += int(containerStatus.RestartCount)
			}
			age := metav1.Now().Sub(pod.CreationTimestamp.Time).Round(0).String()

			table.Append([]string{pod.Name, pod.Namespace, string(pod.Status.Phase), fmt.Sprintf("%d", restarts), age})
		}
		table.Render()
	}
}

// listDeployments function to get and display deployment status
func listDeployments(clientset *kubernetes.Clientset, namespace string) {
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing deployments: %v\n", err)
		return
	}

	fmt.Println("\n--- Deployments ---")
	if len(deployments.Items) == 0 {
		fmt.Println("No deployments found.")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"NAME", "NAMESPACE", "READY", "UP-TO-DATE", "AVAILABLE", "AGE"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")

		for _, deploy := range deployments.Items {
			ready := fmt.Sprintf("%d/%d", deploy.Status.Replicas-deploy.Status.UnavailableReplicas, deploy.Status.Replicas)
			upToDate := fmt.Sprintf("%d", deploy.Status.UpdatedReplicas)
			available := fmt.Sprintf("%d", deploy.Status.AvailableReplicas)
			age := metav1.Now().Sub(deploy.CreationTimestamp.Time).Round(0).String()

			table.Append([]string{deploy.Name, deploy.Namespace, ready, upToDate, available, age})
		}
		table.Render()
	}
}

// listServices function to get and display service status
func listServices(clientset *kubernetes.Clientset, namespace string) {
	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing services: %v\n", err)
		return
	}

	fmt.Println("\n--- Services ---")
	if len(services.Items) == 0 {
		fmt.Println("No services found.")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"NAME", "NAMESPACE", "TYPE", "CLUSTER-IP", "EXTERNAL-IP", "PORTS", "AGE"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")

		for _, svc := range services.Items {
			clusterIP := svc.Spec.ClusterIP
			externalIP := "<none>"
			if len(svc.Status.LoadBalancer.Ingress) > 0 {
				ips := []string{}
				for _, ingress := range svc.Status.LoadBalancer.Ingress {
					if ingress.IP != "" {
						ips = append(ips, ingress.IP)
					} else if ingress.Hostname != "" {
						ips = append(ips, ingress.Hostname)
					}
				}
				if len(ips) > 0 {
					externalIP = strings.Join(ips, ", ")
				}
			}

			ports := []string{}
			for _, port := range svc.Spec.Ports {
				ports = append(ports, fmt.Sprintf("%d:%d/%s", port.Port, port.NodePort, port.Protocol))
			}
			portsStr := strings.Join(ports, ",")
			if portsStr == "" {
				portsStr = "<none>"
			}

			age := metav1.Now().Sub(svc.CreationTimestamp.Time).Round(0).String()

			table.Append([]string{svc.Name, svc.Namespace, string(svc.Spec.Type), clusterIP, externalIP, portsStr, age})
		}
		table.Render()
	}
}

// listIngresses function to get and display ingress status
func listIngresses(clientset *kubernetes.Clientset, namespace string) {
	ingresses, err := clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing ingresses: %v\n", err)
		return
	}

	fmt.Println("\n--- Ingresses ---")
	if len(ingresses.Items) == 0 {
		fmt.Println("No ingresses found.")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"NAME", "NAMESPACE", "CLASS", "HOSTS", "ADDRESS", "PORTS", "AGE"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")

		for _, ing := range ingresses.Items {
			ingressClass := "<none>"
			if ing.Spec.IngressClassName != nil {
				ingressClass = *ing.Spec.IngressClassName
			}

			hosts := []string{}
			for _, rule := range ing.Spec.Rules {
				if rule.Host != "" {
					hosts = append(hosts, rule.Host)
				}
			}
			hostsStr := strings.Join(hosts, ",")
			if hostsStr == "" {
				hostsStr = "*" // Default for ingresses with no specific host rules
			}

			address := "<pending>"
			if len(ing.Status.LoadBalancer.Ingress) > 0 {
				ips := []string{}
				for _, ingStatus := range ing.Status.LoadBalancer.Ingress {
					if ingStatus.IP != "" {
						ips = append(ips, ingStatus.IP)
					} else if ingStatus.Hostname != "" {
						ips = append(ips, ingStatus.Hostname)
					}
				}
				if len(ips) > 0 {
					address = strings.Join(ips, ", ")
				}
			}

			// Common Ingress ports are 80 and 443
			ports := "80"
			for _, tls := range ing.Spec.TLS {
				if len(tls.Hosts) > 0 {
					ports = "80, 443" // If TLS is configured, assume 443 is also open
					break
				}
			}

			age := metav1.Now().Sub(ing.CreationTimestamp.Time).Round(0).String()

			table.Append([]string{ing.Name, ing.Namespace, ingressClass, hostsStr, address, ports, age})
		}
		table.Render()
	}
}
