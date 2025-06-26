package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1" // Importar appsv1
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

var (
	daemonsetsNamespace     string
	daemonsetsAllNamespaces bool
	daemonsetsSelector      string
	daemonsetsOutputFormat  string
)

var daemonsetsGetCmd = &cobra.Command{
	Use:     "daemonsets [nombre-del-daemonset]",
	Aliases: []string{"ds"},
	Short:   "Lista uno o más daemonsets",
	Long:    `Lista uno o más daemonsets en el namespace actual o en todos los namespaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clients, err := GetKubeClients()
		if err != nil {
			return fmt.Errorf("creando clientes de Kubernetes: %w", err)
		}
		listOptions := metav1.ListOptions{LabelSelector: daemonsetsSelector}
		effectiveNamespace := GetEffectiveNamespace(daemonsetsNamespace, daemonsetsAllNamespaces, "default", false)
		if daemonsetsAllNamespaces {
			effectiveNamespace = ""
		}

		var dsList *appsv1.DaemonSetList
		var singleDS *appsv1.DaemonSet

		if len(args) > 0 {
			dsName := args[0]
			currentContextNs := GetEffectiveNamespace(daemonsetsNamespace, false, "default", false)
			if Verbose {
				fmt.Printf("DEBUG: Buscando daemonset '%s' en namespace '%s'\n", dsName, currentContextNs)
			}
			singleDS, err = clients.Core.AppsV1().DaemonSets(currentContextNs).Get(context.TODO(), dsName, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("error obteniendo daemonset '%s': %w", dsName, err)
			}
		} else {
			if Verbose {
				fmt.Printf("DEBUG: Listando daemonsets en namespace '%s' con selector '%s'\n", effectiveNamespace, daemonsetsSelector)
			}
			dsList, err = clients.Core.AppsV1().DaemonSets(effectiveNamespace).List(context.TODO(), listOptions)
			if err != nil {
				return fmt.Errorf("error listando daemonsets: %w", err)
			}
		}

		itemsToProcess := []appsv1.DaemonSet{}
		if singleDS != nil {
			itemsToProcess = append(itemsToProcess, *singleDS)
		} else if dsList != nil {
			itemsToProcess = dsList.Items
		}

		if len(itemsToProcess) == 0 {
			fmt.Fprintln(os.Stdout, "No se encontraron daemonsets.")
			return nil
		}

		outputLower := strings.ToLower(daemonsetsOutputFormat)
		if outputLower == "json" {
			data, _ := json.MarshalIndent(itemsToProcess, "", "  ")
			fmt.Fprintln(os.Stdout, string(data))
			return nil
		}
		if outputLower == "yaml" {
			data, _ := yaml.Marshal(itemsToProcess)
			fmt.Fprintln(os.Stdout, string(data))
			return nil
		}

		headers := []string{"NAMESPACE", "NAME", "DESIRED", "CURRENT", "READY", "UP-TO-DATE", "AVAILABLE", "NODE SELECTOR", "AGE"}
		if outputLower == "wide" {
			headers = append(headers, "CONTAINERS", "IMAGES")
		} // Ejemplo para wide

		rows := make([][]string, 0, len(itemsToProcess))
		for _, ds := range itemsToProcess {
			age := metav1.Now().Sub(ds.CreationTimestamp.Time).Truncate(time.Second).String()
			nodeSelectorStr := metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: ds.Spec.Template.Spec.NodeSelector})
			if nodeSelectorStr == "" {
				nodeSelectorStr = "<none>"
			}

			containers := []string{}
			images := []string{}
			for _, c := range ds.Spec.Template.Spec.Containers {
				containers = append(containers, c.Name)
				images = append(images, c.Image)
			}

			row := []string{
				ds.Namespace, ds.Name,
				fmt.Sprintf("%d", ds.Status.DesiredNumberScheduled),
				fmt.Sprintf("%d", ds.Status.CurrentNumberScheduled),
				fmt.Sprintf("%d", ds.Status.NumberReady),
				fmt.Sprintf("%d", ds.Status.UpdatedNumberScheduled),
				fmt.Sprintf("%d", ds.Status.NumberAvailable),
				nodeSelectorStr, age,
			}
			if outputLower == "wide" {
				row = append(row, strings.Join(containers, ","), strings.Join(images, ","))
			}
			rows = append(rows, row)
		}
		PrintBasicTable(headers, rows)
		return nil
	},
}

func init() {
	getCmd.AddCommand(daemonsetsGetCmd)
	daemonsetsGetCmd.Flags().StringVarP(&daemonsetsNamespace, "namespace", "n", "", "Namespace")
	daemonsetsGetCmd.Flags().BoolVarP(&daemonsetsAllNamespaces, "all-namespaces", "A", false, "Todos los namespaces")
	daemonsetsGetCmd.Flags().StringVarP(&daemonsetsSelector, "selector", "l", "", "Selector (label query)")
	daemonsetsGetCmd.Flags().StringVarP(&daemonsetsOutputFormat, "output", "o", "", "Formato de salida (wide, json, yaml)")
}
