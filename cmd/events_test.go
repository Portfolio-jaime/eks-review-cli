// cmd/events_test.go
package cmd

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubefake "k8s.io/client-go/kubernetes/fake"
	// ktesting "k8s.io/client-go/testing" // Para Reactors más avanzados
)

// Helper para ejecutar comandos de Cobra y capturar salida
func executeCommand(root *cobra.Command, args ...string) (stdout, stderr string, err error) {
	stdoutBuf := new(bytes.Buffer)
	stderrBuf := new(bytes.Buffer)
	root.SetOut(stdoutBuf)
	root.SetErr(stderrBuf)
	root.SetArgs(args)

	err = root.Execute() // ExecuteC podría ser útil si necesitas el comando ejecutado

	return stdoutBuf.String(), stderrBuf.String(), err
}

func TestEventsCmd(t *testing.T) {
	// Guardar y restaurar la función original GetKubeClients
	// y la variable Verbose para aislar la prueba.
	originalGetKubeClients := GetKubeClients // Asumiendo que es una variable global o accesible
	originalVerbose := Verbose
	defer func() {
		GetKubeClients = originalGetKubeClients
		Verbose = originalVerbose
	}()

	// Mock de eventos
	mockEventTime := metav1.Now()
	mockEvents := []runtime.Object{
		&corev1.Event{
			ObjectMeta:     metav1.ObjectMeta{Name: "event1", Namespace: "default"},
			InvolvedObject: corev1.ObjectReference{Kind: "Pod", Name: "test-pod"},
			Reason:         "Scheduled", Type: "Normal", Message: "Pod successfully scheduled",
			LastTimestamp: mockEventTime,
		},
		&corev1.Event{
			ObjectMeta:     metav1.ObjectMeta{Name: "event2", Namespace: "default"},
			InvolvedObject: corev1.ObjectReference{Kind: "Node", Name: "node1"},
			Reason:         "NodeNotReady", Type: "Warning", Message: "Node is not ready",
			LastTimestamp: metav1.NewTime(mockEventTime.Add(-5 * time.Minute)),
		},
	}

	fakeClientset := kubefake.NewSimpleClientset(mockEvents...)

	// Sobrescribir GetKubeClients para que devuelva nuestro fake clientset
	GetKubeClients = func() *KubeClients {
		return &KubeClients{Core: fakeClientset, Metrics: nil} // Asumimos que no se necesitan métricas para 'events'
	}
	Verbose = false // Deshabilitar DEBUG logs para la prueba de salida estándar

	// Crear una instancia del comando raíz para probar (o el comando monitor directamente si es más fácil)
	// Asumimos que rootCmd está accesible o tienes una forma de obtener eventsCmd
	// Esto es más fácil si `rootCmd` y sus subcomandos se configuran en una función que puedas llamar.
	// Por ahora, usaré directamente `eventsCmd` si está exportado o accesible.
	// Si no, tendrías que ejecutar a través de `rootCmd`.

	// Necesitas resetear las flags para cada ejecución de prueba si son persistentes
	// o usar cmd.ResetFlags() si el comando lo permite.

	t.Run("Muestra eventos normalmente", func(t *testing.T) {
		// Reiniciar flags si es necesario, o asegurarse de que estén en su estado por defecto
		eventType = ""       // Valor por defecto
		eventsNamespace = "" // Valor por defecto (GetEffectiveNamespace se encargará del default "default")

		// Ejecutar `eks-review-cli monitor events` (simulado)
		// Esto requiere que `rootCmd` esté configurado y `monitorCmd` y `eventsCmd` añadidos.
		// Para simplificar, si `eventsCmd.Run` es tu lógica principal, puedes llamarla directamente,
		// pero es mejor probar a través de la ejecución del comando.
		// Aquí, asumo que podemos ejecutar el subcomando events directamente para probar su salida.

		// Re-crear el árbol de comandos para esta prueba para asegurar estado limpio
		testRootCmd, _, _ := GetTestRootCommand() // Necesitarías una función helper para esto
		stdout, stderr, err := executeCommandC(testRootCmd, "monitor", "events")

		require.NoError(t, err)
		assert.Empty(t, stderr, "Stderr debería estar vacío")

		// Verificar la salida (esto es frágil, pero necesario para probar tablas)
		assert.Contains(t, stdout, "event1")
		assert.Contains(t, stdout, "test-pod")
		assert.Contains(t, stdout, "Normal")
		assert.Contains(t, stdout, "event2")
		assert.Contains(t, stdout, "NodeNotReady")
		assert.Contains(t, stdout, "Warning")
		assert.Contains(t, stdout, "default") // Namespace
		// Podrías hacer aserciones más específicas sobre el formato de la tabla
	})

	t.Run("Filtra eventos por tipo Warning", func(t *testing.T) {
		testRootCmd, _, _ := GetTestRootCommand()
		stdout, stderr, err := executeCommandC(testRootCmd, "monitor", "events", "--type", "Warning")

		require.NoError(t, err)
		assert.Empty(t, stderr)
		assert.NotContains(t, stdout, "event1") // Normal event
		assert.Contains(t, stdout, "event2")    // Warning event
		assert.Contains(t, stdout, "NodeNotReady")
	})

	t.Run("Filtra por namespace (si GetEffectiveNamespace funciona con el fake kubeconfig)", func(t *testing.T) {
		// Este caso es más complejo porque GetEffectiveNamespace usa el filesystem.
		// Para probarlo bien, necesitarías que GetEffectiveNamespace también sea mockeable
		// o que el KUBECONFIG del entorno de prueba apunte a un kubeconfig falso
		// que defina un namespace diferente para el contexto actual.

		// Aquí, vamos a simular que el namespace se pasa directamente como flag
		// y que GetEffectiveNamespace lo usa.
		fakeClientsetWithNs := kubefake.NewSimpleClientset(
			&corev1.Event{
				ObjectMeta:     metav1.ObjectMeta{Name: "event-ns", Namespace: "custom-ns"},
				InvolvedObject: corev1.ObjectReference{Kind: "Pod", Name: "pod-ns"},
				Reason:         "Pulled", Type: "Normal", Message: "Image pulled",
				LastTimestamp: mockEventTime,
			},
		)
		GetKubeClients = func() *KubeClients {
			return &KubeClients{Core: fakeClientsetWithNs}
		}

		testRootCmd, _, _ := GetTestRootCommand()
		// Resetear la flag global eventsNamespace antes de establecerla
		// o, mejor, buscar cómo Cobra maneja el reseteo de flags entre pruebas.
		// Para pruebas, a veces es más fácil instanciar el comando de nuevo.
		stdout, stderr, err := executeCommandC(testRootCmd, "monitor", "events", "-n", "custom-ns")

		require.NoError(t, err)
		assert.Empty(t, stderr)
		assert.Contains(t, stdout, "event-ns")
		assert.Contains(t, stdout, "custom-ns")
		assert.NotContains(t, stdout, "default") // No debería mostrar eventos del namespace 'default'
	})

	// TODO: Añadir pruebas para --verbose, errores de API (usando Reactors), etc.
}

// executeCommandC es una versión de executeCommand que devuelve el comando ejecutado.
// Necesitarás una función para obtener tu rootCmd configurado para las pruebas.
// Esta es una simplificación.
func GetTestRootCommand() (rootCmd *cobra.Command, monitorCmd *cobra.Command, eventsCmd *cobra.Command) {
	// Esto recrearía la estructura de tu comando para la prueba
	// Necesitarías re-instanciar tus comandos aquí para evitar estado contaminado entre pruebas.
	// Esta es una parte crucial y a veces compleja de las pruebas de Cobra.
	// La estructura sería similar a tu main.go y cmd/root.go, cmd/monitor.go
	r := &cobra.Command{Use: "eks-review-cli-test"}
	r.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose output") // Asegúrate que la flag Verbose esté aquí

	m := &cobra.Command{Use: "monitor"}

	// Re-declarar tus variables de flags locales para eventsCmd si no son globales
	// var testEventType string
	// var testEventsNamespace string

	eCmd := &cobra.Command{
		Use: "events",
		Run: func(cmd *cobra.Command, args []string) {
			// Aquí copiarías la lógica de tu eventsCmd.Run original,
			// o mejor, tu eventsCmd.RunE si devuelve un error.
			// Esto se vuelve complicado si las flags están definidas como variables globales en el paquete cmd.
			// Es una de las razones por las que a veces se prefiere pasar dependencias (como los valores de las flags)
			// directamente a las funciones de lógica.

			// Para este ejemplo, asumimos que la lógica de eventsCmd.Run accede a las flags globales del paquete.
			// Es importante que esas flags se reseteen antes de cada sub-prueba si es necesario.
			// Por simplicidad, el `eventsCmd` real se usa en `executeCommand`,
			// pero debes tener cuidado con el estado de las flags.
			// Aquí solo devolvemos la estructura.
		},
	}
	// Adjuntar flags a eCmd como lo haces en eventsCmd.init()
	eCmd.Flags().StringVarP(&eventType, "type", "T", "", "Filter events by type")
	eCmd.Flags().StringVarP(&eventsNamespace, "namespace", "n", "", "Namespace scope")

	m.AddCommand(eCmd)
	r.AddCommand(m)
	return r, m, eCmd
}

// executeCommandC para obtener el comando y sus buffers
func executeCommandC(root *cobra.Command, args ...string) (stdout string, stderr string, err error) {
	stdoutBuf := new(bytes.Buffer)
	stderrBuf := new(bytes.Buffer)

	actualCmd, err := findSubcommand(root, args...)
	if err != nil {
		return "", "", fmt.Errorf("subcommand not found for args %v: %w", args, err)
	}

	// Asegurar que Out y Err se establezcan en el comando más específico que se va a ejecutar
	// o en el root si se propagan. Cobra maneja esto, pero para capturar,
	// es bueno establecerlo en el comando que realmente se ejecuta o su raíz.
	root.SetOut(stdoutBuf) // Captura en el root
	root.SetErr(stderrBuf)

	root.SetArgs(args) // Establece los args en el root

	cmdErr := root.Execute()

	return stdoutBuf.String(), stderrBuf.String(), cmdErr
}

// findSubcommand es un helper para obtener el comando específico que se ejecutará.
func findSubcommand(rootCmd *cobra.Command, args ...string) (*cobra.Command, error) {
	cmd, _, err := rootCmd.Find(args)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// Nota: Probar los comandos de Cobra de esta manera puede ser un poco verboso
// debido a la necesidad de configurar el comando y sus flags correctamente para cada caso de prueba.
// Considera refactorizar la lógica principal de tus `Run` functions en métodos o funciones separadas
// que acepten el `clientset` y los valores de las flags como parámetros,
// para que puedas probar esa lógica de forma más aislada y luego tener pruebas de comando más ligeras
// que solo verifiquen el parseo de flags y la llamada a esa lógica.
