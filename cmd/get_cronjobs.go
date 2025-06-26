package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	batchv1 "k8s.io/api/batch/v1" // Para CronJobs (batchv1)
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

var (
	cronjobsNamespace     string
	cronjobsAllNamespaces bool
	cronjobsSelector      string // CronJobs no suelen tener selectores de etiquetas propios, sino que su jobTemplate sí
	cronjobsOutputFormat  string
)

var cronjobsGetCmd = &cobra.Command{
	Use:     "cronjobs [nombre-del-cronjob]",
	Aliases: []string{"cj"},
	Short:   "Lista uno o más cronjobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		clients, err := GetKubeClients()
		if err != nil {
			return fmt.Errorf("creando clientes de Kubernetes: %w", err)
		}
		// CronJobs no son filtrables por label selector directamente a nivel de lista de CronJob,
		// pero lo mantenemos por consistencia si se implementa un filtro custom.
		listOptions := metav1.ListOptions{LabelSelector: cronjobsSelector}
		effectiveNamespace := GetEffectiveNamespace(cronjobsNamespace, cronjobsAllNamespaces, "default", false)
		if cronjobsAllNamespaces {
			effectiveNamespace = ""
		}

		var cjList *batchv1.CronJobList
		var singleCJ *batchv1.CronJob

		if len(args) > 0 {
			cjName := args[0]
			currentContextNs := GetEffectiveNamespace(cronjobsNamespace, false, "default", false)
			singleCJ, err = clients.Core.BatchV1().CronJobs(currentContextNs).Get(context.TODO(), cjName, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("error obteniendo cronjob '%s': %w", cjName, err)
			}
		} else {
			cjList, err = clients.Core.BatchV1().CronJobs(effectiveNamespace).List(context.TODO(), listOptions)
			if err != nil {
				return fmt.Errorf("error listando cronjobs: %w", err)
			}
		}

		itemsToProcess := []batchv1.CronJob{}
		if singleCJ != nil {
			itemsToProcess = append(itemsToProcess, *singleCJ)
		} else if cjList != nil {
			itemsToProcess = cjList.Items
		}

		if len(itemsToProcess) == 0 {
			fmt.Fprintln(os.Stdout, "No se encontraron cronjobs.")
			return nil
		}

		outputLower := strings.ToLower(cronjobsOutputFormat)
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

		headers := []string{"NAMESPACE", "NAME", "SCHEDULE", "SUSPEND", "ACTIVE", "LAST SCHEDULE", "AGE"}
		if outputLower == "wide" {
			headers = append(headers, "LAST SUCCESSFUL TIME")
		}

		rows := make([][]string, 0, len(itemsToProcess))
		for _, cj := range itemsToProcess {
			suspend := "False"
			if cj.Spec.Suspend != nil && *cj.Spec.Suspend {
				suspend = "True"
			}
			lastSchedule := "<none>"
			if cj.Status.LastScheduleTime != nil {
				lastSchedule = metav1.Now().Sub(cj.Status.LastScheduleTime.Time).Truncate(time.Second).String() + " ago"
			}
			age := metav1.Now().Sub(cj.CreationTimestamp.Time).Truncate(time.Second).String()

			lastSuccessfulTimeStr := "<none>"
			if cj.Status.LastSuccessfulTime != nil {
				lastSuccessfulTimeStr = cj.Status.LastSuccessfulTime.Format(time.RFC3339)
			}

			row := []string{
				cj.Namespace, cj.Name, cj.Spec.Schedule, suspend,
				fmt.Sprintf("%d", len(cj.Status.Active)),
				lastSchedule, age,
			}
			if outputLower == "wide" {
				row = append(row, lastSuccessfulTimeStr)
			}
			rows = append(rows, row)
		}
		PrintBasicTable(headers, rows)
		return nil
	},
}

func init() {
	getCmd.AddCommand(cronjobsGetCmd)
	cronjobsGetCmd.Flags().StringVarP(&cronjobsNamespace, "namespace", "n", "", "Namespace")
	cronjobsGetCmd.Flags().BoolVarP(&cronjobsAllNamespaces, "all-namespaces", "A", false, "Todos los namespaces")
	cronjobsGetCmd.Flags().StringVarP(&cronjobsSelector, "selector", "l", "", "Selector (label query) - aplica al job template")
	cronjobsGetCmd.Flags().StringVarP(&cronjobsOutputFormat, "output", "o", "", "Formato de salida (wide, json, yaml)")
}
