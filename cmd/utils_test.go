// cmd/utils_test.go
package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEffectiveNamespace(t *testing.T) {
	// Crear un kubeconfig temporal para pruebas
	tempDir := t.TempDir()
	dummyKubeconfigContent := `
apiVersion: v1
clusters:
- cluster:
    server: https://localhost:8080
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
    namespace: context-namespace
  name: test-context
current-context: test-context
kind: Config
preferences: {}
users:
- name: test-user
`
	dummyKubeconfigPath := filepath.Join(tempDir, "config")
	err := os.WriteFile(dummyKubeconfigPath, []byte(dummyKubeconfigContent), 0600)
	assert.NoError(t, err)

	// Sobrescribir temporalmente la función que obtiene el homedir para apuntar a nuestro kubeconfig
	// O, mejor aún, modificar GetEffectiveNamespace para aceptar el path del kubeconfig directamente.
	// La versión que te di de GetEffectiveNamespace ya usa un path fijo o el default,
	// para probar la parte del contexto actual, necesitarías simular el kubeconfig en esa ubicación.
	// Aquí simplificaré asumiendo que podemos controlar el path o el contexto actual de alguna manera (puede requerir mocks más avanzados o refactorización).

	// Para simplificar esta prueba, nos enfocaremos en la lógica de las flags y el default.
	// Para probar la lógica del contexto actual, necesitarías asegurar que `clientcmd.NewNonInteractiveDeferredLoadingClientConfig`
	// pueda ser influenciado en el entorno de prueba (ej. estableciendo KUBECONFIG env var o mockeando homedir).

	testCases := []struct {
		name                   string
		namespaceFlag          string
		allNamespacesFlag      bool
		defaultNS              string
		commandAllowsAllString bool
		// mockKubeconfigPath string // Si modificas GetEffectiveNamespace para aceptar path
		expectedNamespace string
	}{
		{
			name:              "Flag específica provista",
			namespaceFlag:     "my-ns",
			allNamespacesFlag: false,
			defaultNS:         "default",
			expectedNamespace: "my-ns",
		},
		{
			name:              "AllNamespaces flag es true",
			namespaceFlag:     "ignored-ns",
			allNamespacesFlag: true,
			defaultNS:         "default",
			expectedNamespace: "", // Espera "" para todos los namespaces
		},
		{
			name:                   "Namespace flag es 'all' y el comando lo permite",
			namespaceFlag:          "all",
			allNamespacesFlag:      false,
			defaultNS:              "default",
			commandAllowsAllString: true,
			expectedNamespace:      "",
		},
		{
			name:                   "Namespace flag es 'All' (case-insensitive) y el comando lo permite",
			namespaceFlag:          "All",
			allNamespacesFlag:      false,
			defaultNS:              "default",
			commandAllowsAllString: true,
			expectedNamespace:      "",
		},
		{
			name:                   "Namespace flag es 'other' y el comando permite 'all' (no debe ser 'all')",
			namespaceFlag:          "other",
			allNamespacesFlag:      false,
			defaultNS:              "default",
			commandAllowsAllString: true,
			expectedNamespace:      "other",
		},
		{
			name:              "Sin flags, usa default",
			namespaceFlag:     "",
			allNamespacesFlag: false,
			defaultNS:         "test-default",
			expectedNamespace: "test-default",
		},
		// TODO: Añadir casos para probar la lectura del contexto actual del kubeconfig.
		// Esto requeriría configurar el entorno KUBECONFIG para que apunte a dummyKubeconfigPath
		// o modificar GetEffectiveNamespace para que acepte el path del kubeconfig.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Si tu función GetEffectiveNamespace depende de homedir.HomeDir()
			// y no puedes mockearlo fácilmente, esta parte de la prueba del contexto actual
			// será más difícil de aislar.
			// La versión de GetEffectiveNamespace que te dí intenta cargar el kubeconfig desde el homedir.
			// Para probar esa parte específica, podrías:
			// 1. Establecer la variable de entorno KUBECONFIG al 'dummyKubeconfigPath' ANTES de llamar a la función.
			//    originalKubeconfig := os.Getenv("KUBECONFIG")
			//    os.Setenv("KUBECONFIG", dummyKubeconfigPath)
			//    defer os.Setenv("KUBECONFIG", originalKubeconfig) // Restaurar
			//
			//    O modificar GetEffectiveNamespace para que use clientcmd.NewDefaultClientConfigLoadingRules()
			//    y puedas establecer ExplicitPath en las reglas de carga para la prueba.

			actualNamespace := GetEffectiveNamespace(tc.namespaceFlag, tc.allNamespacesFlag, tc.defaultNS, tc.commandAllowsAllString)
			assert.Equal(t, tc.expectedNamespace, actualNamespace)
		})
	}
}
