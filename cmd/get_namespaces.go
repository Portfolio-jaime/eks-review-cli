// cmd/get_namespaces.go
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1" // Para Namespaces
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

var (
	// No necesitamos flags de namespace para listar namespaces
	namespacesSelector     string
	namespacesOutputFormat string
)

var namespacesGetCmd = &cobra.Command{
	Use:     "namespaces [nombre-del-namespace]",
	Aliases: []string{"ns"},
	Short:   "Lista uno o más namespaces",
	RunE: func(cmd *cobra.Command, args []string) error {
		clients := GetKubeClients()
		listOptions := metav1.ListOptions{LabelSelector: namespacesSelector}

		var nsList *corev1.NamespaceList
		var singleNS *corev1.Namespace
		var err error

		if len(args) > 0 {
			nsName := args[0]
			singleNS, err = clients.Core.CoreV1().Namespaces().Get(context.TODO(), nsName, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("error obteniendo namespace '%s': %w", nsName, err)
			}
		} else {
			nsList, err = clients.Core.CoreV1().Namespaces().List(context.TODO(), listOptions)
			if err != nil {
				return fmt.Errorf("error listando namespaces: %w", err)
			}
		}

		itemsToProcess := []corev1.Namespace{}
		if singleNS != nil {
			itemsToProcess = append(itemsToProcess, *singleNS)
		} else if nsList != nil {
			itemsToProcess = nsList.Items
		}

		if len(itemsToProcess) == 0 {
			fmt.Fprintln(os.Stdout, "No se encontraron namespaces.")
			return nil
		}

		outputLower := strings.ToLower(namespacesOutputFormat)
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

		headers := []string{"NAME", "STATUS", "AGE"}
		// Para -o wide podríamos añadir LABELS, pero sería mucho texto
		// if outputLower == "wide" { headers = append(headers, "LABELS") }

		rows := make([][]string, 0, len(itemsToProcess))
		for _, ns := range itemsToProcess {
			age := metav1.Now().Sub(ns.CreationTimestamp.Time).Truncate(time.Second).String()
			row := []string{ns.Name, string(ns.Status.Phase), age}
			// if outputLower == "wide" { row = append(row, metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: ns.Labels}))}
			rows = append(rows, row)
		}
		PrintBasicTable(headers, rows)
		return nil
	},
}

func init() {
	getCmd.AddCommand(namespacesGetCmd)
	// No flags de namespace aquí
	namespacesGetCmd.Flags().StringVarP(&namespacesSelector, "selector", "l", "", "Selector (label query) para filtrar namespaces")
	namespacesGetCmd.Flags().StringVarP(&namespacesOutputFormat, "output", "o", "", "Formato de salida (wide, json, yaml)")
}
