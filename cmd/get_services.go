package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

var (
	servicesNamespace     string
	servicesAllNamespaces bool
	servicesSelector      string
	servicesOutputFormat  string
)

var servicesGetCmd = &cobra.Command{
	Use:     "services [nombre-del-servicio]",
	Aliases: []string{"svc"},
	Short:   "Lista uno o más services",
	Long:    `Lista uno o más services en el namespace actual o en todos los namespaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clients := GetKubeClients()
		listOptions := metav1.ListOptions{LabelSelector: servicesSelector}
		effectiveNamespace := GetEffectiveNamespace(servicesNamespace, servicesAllNamespaces, "default", false)
		if servicesAllNamespaces {
			effectiveNamespace = ""
		}

		var serviceList *corev1.ServiceList
		var singleService *corev1.Service
		var err error

		if len(args) > 0 {
			serviceName := args[0]
			currentContextNs := GetEffectiveNamespace(servicesNamespace, false, "default", false)
			if Verbose {
				fmt.Printf("DEBUG: Buscando service '%s' en namespace '%s'\n", serviceName, currentContextNs)
			}
			singleService, err = clients.Core.CoreV1().Services(currentContextNs).Get(context.TODO(), serviceName, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("error obteniendo service '%s': %w", serviceName, err)
			}
		} else {
			if Verbose {
				fmt.Printf("DEBUG: Listando services en namespace '%s' con selector '%s'\n", effectiveNamespace, servicesSelector)
			}
			serviceList, err = clients.Core.CoreV1().Services(effectiveNamespace).List(context.TODO(), listOptions)
			if err != nil {
				return fmt.Errorf("error listando services: %w", err)
			}
		}

		itemsToProcess := []corev1.Service{}
		if singleService != nil {
			itemsToProcess = append(itemsToProcess, *singleService)
		} else if serviceList != nil {
			itemsToProcess = serviceList.Items
		}

		if len(itemsToProcess) == 0 {
			fmt.Fprintln(os.Stdout, "No se encontraron services.")
			return nil
		}

		outputLower := strings.ToLower(servicesOutputFormat)
		if outputLower == "json" { /* ... (igual que en pods) ... */
			data, _ := json.MarshalIndent(itemsToProcess, "", "  ")
			fmt.Fprintln(os.Stdout, string(data))
			return nil
		}
		if outputLower == "yaml" { /* ... (igual que en pods) ... */
			data, _ := yaml.Marshal(itemsToProcess)
			fmt.Fprintln(os.Stdout, string(data))
			return nil
		}

		headers := []string{"NAMESPACE", "NAME", "TYPE", "CLUSTER-IP", "EXTERNAL-IP", "PORT(S)", "AGE"}
		if outputLower == "wide" {
			headers = append(headers, "SELECTOR")
		}

		rows := make([][]string, 0, len(itemsToProcess))
		for _, svc := range itemsToProcess {
			externalIP := "<none>"
			if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
				if len(svc.Status.LoadBalancer.Ingress) > 0 {
					ips := []string{}
					for _, ing := range svc.Status.LoadBalancer.Ingress {
						if ing.IP != "" {
							ips = append(ips, ing.IP)
						}
						if ing.Hostname != "" {
							ips = append(ips, ing.Hostname)
						}
					}
					if len(ips) > 0 {
						externalIP = strings.Join(ips, ",")
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

			selectorStr := metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: svc.Spec.Selector})

			row := []string{svc.Namespace, svc.Name, string(svc.Spec.Type), svc.Spec.ClusterIP, externalIP, strings.Join(portStrings, ","), age}
			if outputLower == "wide" {
				row = append(row, selectorStr)
			}
			rows = append(rows, row)
		}
		PrintBasicTable(headers, rows)
		return nil
	},
}

func init() {
	getCmd.AddCommand(servicesGetCmd)
	servicesGetCmd.Flags().StringVarP(&servicesNamespace, "namespace", "n", "", "Namespace")
	servicesGetCmd.Flags().BoolVarP(&servicesAllNamespaces, "all-namespaces", "A", false, "Todos los namespaces")
	servicesGetCmd.Flags().StringVarP(&servicesSelector, "selector", "l", "", "Selector (label query)")
	servicesGetCmd.Flags().StringVarP(&servicesOutputFormat, "output", "o", "", "Formato de salida (wide, json, yaml)")
}
