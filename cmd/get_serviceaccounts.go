// cmd/get_serviceaccounts.go
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1" // Para ServiceAccounts
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

var (
	saNamespace     string
	saAllNamespaces bool
	saSelector      string // ServiceAccounts pueden tener labels
	saOutputFormat  string
)

var serviceaccountsGetCmd = &cobra.Command{
	Use:     "serviceaccounts [nombre-del-sa]",
	Aliases: []string{"sa"},
	Short:   "Lista uno o más serviceaccounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		clients := GetKubeClients()
		listOptions := metav1.ListOptions{LabelSelector: saSelector}
		effectiveNamespace := GetEffectiveNamespace(saNamespace, saAllNamespaces, "default", false)
		if saAllNamespaces {
			effectiveNamespace = ""
		}

		var saList *corev1.ServiceAccountList
		var singleSA *corev1.ServiceAccount
		var err error

		if len(args) > 0 {
			saName := args[0]
			currentContextNs := GetEffectiveNamespace(saNamespace, false, "default", false)
			singleSA, err = clients.Core.CoreV1().ServiceAccounts(currentContextNs).Get(context.TODO(), saName, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("error obteniendo serviceaccount '%s': %w", saName, err)
			}
		} else {
			saList, err = clients.Core.CoreV1().ServiceAccounts(effectiveNamespace).List(context.TODO(), listOptions)
			if err != nil {
				return fmt.Errorf("error listando serviceaccounts: %w", err)
			}
		}

		itemsToProcess := []corev1.ServiceAccount{}
		if singleSA != nil {
			itemsToProcess = append(itemsToProcess, *singleSA)
		} else if saList != nil {
			itemsToProcess = saList.Items
		}

		if len(itemsToProcess) == 0 {
			fmt.Fprintln(os.Stdout, "No se encontraron serviceaccounts.")
			return nil
		}

		outputLower := strings.ToLower(saOutputFormat)
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

		headers := []string{"NAMESPACE", "NAME", "SECRETS", "AGE"}
		if outputLower == "wide" {
			headers = append(headers, "AUTOMOUNT")
		}

		rows := make([][]string, 0, len(itemsToProcess))
		for _, sa := range itemsToProcess {
			age := metav1.Now().Sub(sa.CreationTimestamp.Time).Truncate(time.Second).String()
			secretsCount := len(sa.Secrets) // Número de secrets referenciados (puede no ser lo mismo que montados)

			automount := "<nil>"
			if sa.AutomountServiceAccountToken != nil {
				automount = fmt.Sprintf("%t", *sa.AutomountServiceAccountToken)
			}

			row := []string{sa.Namespace, sa.Name, fmt.Sprintf("%d", secretsCount), age}
			if outputLower == "wide" {
				row = append(row, automount)
			}
			rows = append(rows, row)
		}
		PrintBasicTable(headers, rows)
		return nil
	},
}

func init() {
	getCmd.AddCommand(serviceaccountsGetCmd)
	serviceaccountsGetCmd.Flags().StringVarP(&saNamespace, "namespace", "n", "", "Namespace")
	serviceaccountsGetCmd.Flags().BoolVarP(&saAllNamespaces, "all-namespaces", "A", false, "Todos los namespaces")
	serviceaccountsGetCmd.Flags().StringVarP(&saSelector, "selector", "l", "", "Selector (label query)")
	serviceaccountsGetCmd.Flags().StringVarP(&saOutputFormat, "output", "o", "", "Formato de salida (wide, json, yaml)")
}
