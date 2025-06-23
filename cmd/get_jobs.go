package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	batchv1 "k8s.io/api/batch/v1" // Para Jobs
	corev1 "k8s.io/api/core/v1"   // Necesario para corev1.ConditionTrue
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

var (
	jobsNamespace     string
	jobsAllNamespaces bool
	jobsSelector      string
	jobsOutputFormat  string
)

var jobsGetCmd = &cobra.Command{
	Use:     "jobs [nombre-del-job]",
	Aliases: []string{"job"},
	Short:   "Lista uno o más jobs",
	Long:    `Lista uno o más jobs en el namespace actual o en todos los namespaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clients := GetKubeClients()
		listOptions := metav1.ListOptions{LabelSelector: jobsSelector}
		effectiveNamespace := GetEffectiveNamespace(jobsNamespace, jobsAllNamespaces, "default", false)
		if jobsAllNamespaces {
			effectiveNamespace = ""
		}

		var jobList *batchv1.JobList
		var singleJob *batchv1.Job
		var err error

		if len(args) > 0 { // Si se especifica un nombre de job
			jobName := args[0]
			currentContextNs := GetEffectiveNamespace(jobsNamespace, false, "default", false) // Namespace para un get individual
			if Verbose {
				fmt.Printf("DEBUG: Buscando job '%s' en namespace '%s'\n", jobName, currentContextNs)
			}
			singleJob, err = clients.Core.BatchV1().Jobs(currentContextNs).Get(context.TODO(), jobName, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("error obteniendo job '%s' en namespace '%s': %w", jobName, currentContextNs, err)
			}
		} else { // Si no se especifica nombre de job, listar jobs
			if Verbose {
				fmt.Printf("DEBUG: Listando jobs en namespace '%s' con selector '%s'\n", effectiveNamespace, jobsSelector)
			}
			jobList, err = clients.Core.BatchV1().Jobs(effectiveNamespace).List(context.TODO(), listOptions)
			if err != nil {
				return fmt.Errorf("error listando jobs: %w", err)
			}
		}

		itemsToProcess := []batchv1.Job{}
		if singleJob != nil {
			itemsToProcess = append(itemsToProcess, *singleJob)
		} else if jobList != nil {
			itemsToProcess = jobList.Items
		}

		if len(itemsToProcess) == 0 {
			fmt.Fprintln(os.Stdout, "No se encontraron jobs.")
			return nil
		}

		outputLower := strings.ToLower(jobsOutputFormat)
		if outputLower == "json" {
			data, errJson := json.MarshalIndent(itemsToProcess, "", "  ")
			if errJson != nil {
				return fmt.Errorf("error convirtiendo a JSON: %w", errJson)
			}
			fmt.Fprintln(os.Stdout, string(data))
			return nil
		} else if outputLower == "yaml" {
			data, errYaml := yaml.Marshal(itemsToProcess)
			if errYaml != nil {
				return fmt.Errorf("error convirtiendo a YAML: %w", errYaml)
			}
			fmt.Fprintln(os.Stdout, string(data))
			return nil
		}

		headers := []string{"NAMESPACE", "NAME", "COMPLETIONS", "DURATION", "AGE"}
		if outputLower == "wide" {
			headers = append(headers, "CONDITIONS", "SELECTOR")
		}

		rows := make([][]string, 0, len(itemsToProcess))
		for _, job := range itemsToProcess {
			completions := "N/A"
			if job.Spec.Completions != nil {
				completions = fmt.Sprintf("%d/%d", job.Status.Succeeded, *job.Spec.Completions)
			}

			duration := "<none>"
			if job.Status.StartTime != nil && job.Status.CompletionTime != nil {
				duration = job.Status.CompletionTime.Sub(job.Status.StartTime.Time).Truncate(time.Second).String()
			} else if job.Status.StartTime != nil {
				// Job todavía corriendo o no ha completado/fallado con tiempo de finalización
				duration = time.Since(job.Status.StartTime.Time).Truncate(time.Second).String()
				// Se podría añadir un sufijo como "(running)" si Status.Active > 0
				if job.Status.Active > 0 {
					duration += " (running)"
				}
			}

			age := metav1.Now().Sub(job.CreationTimestamp.Time).Truncate(time.Second).String()

			conditions := []string{}
			for _, cond := range job.Status.Conditions {
				// CORRECCIÓN AQUÍ: Usar corev1.ConditionTrue
				if cond.Status == corev1.ConditionTrue {
					conditions = append(conditions, string(cond.Type))
				}
			}
			conditionsStr := strings.Join(conditions, ",")
			if conditionsStr == "" {
				conditionsStr = "<none>"
			}

			selectorStr := "<none>"
			if job.Spec.Selector != nil {
				selectorStr = metav1.FormatLabelSelector(job.Spec.Selector)
			}

			row := []string{job.Namespace, job.Name, completions, duration, age}
			if outputLower == "wide" {
				row = append(row, conditionsStr, selectorStr)
			}
			rows = append(rows, row)
		}
		PrintBasicTable(headers, rows)
		return nil
	},
}

func init() {
	getCmd.AddCommand(jobsGetCmd)
	jobsGetCmd.Flags().StringVarP(&jobsNamespace, "namespace", "n", "", "Namespace para listar jobs (opcional)")
	jobsGetCmd.Flags().BoolVarP(&jobsAllNamespaces, "all-namespaces", "A", false, "Listar jobs en todos los namespaces")
	jobsGetCmd.Flags().StringVarP(&jobsSelector, "selector", "l", "", "Selector (label query) para filtrar jobs")
	jobsGetCmd.Flags().StringVarP(&jobsOutputFormat, "output", "o", "", "Formato de salida. Soportado: wide, json, yaml")
}
