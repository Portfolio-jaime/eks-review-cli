// cmd/logs_test.go
package cmd

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLineScanner(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedLines []string
		expectedErr   error
	}{
		{"Vacío", "", []string{}, io.EOF}, // Esperar io.EOF si no hay nada que escanear la primera vez
		{"Una línea", "hello\n", []string{"hello"}, nil},
		{"Múltiples líneas", "hello\nworld\n", []string{"hello", "world"}, nil},
		{"Sin newline al final", "hello\nworld", []string{"hello", "world"}, nil},
		{"Solo newline", "\n", []string{""}, nil},
		{"Múltiples newlines", "\n\nhello\n", []string{"", "", "hello"}, nil},
		// Podrías añadir casos con buffers pequeños si quieres probar la expansión del buffer,
		// pero eso es más complejo de simular directamente sin acceder a internos.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.input)
			scanner := NewLineScanner(reader)
			var lines []string
			var err error

			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			err = scanner.Err()

			// Si la entrada está vacía y no se llama a Scan(), Err() será nil.
			// Si se llama a Scan() sobre una entrada vacía, Scan() devuelve false y Err() será io.EOF.
			if tc.input == "" {
				assert.Equal(t, 0, len(lines), "No debería haber líneas para entrada vacía después de escanear")
				require.ErrorIs(t, err, io.EOF, "Err() debería ser io.EOF para entrada vacía después de escanear")
			} else {
				assert.Equal(t, tc.expectedLines, lines)
				if tc.expectedErr != nil {
					require.ErrorIs(t, err, tc.expectedErr)
				} else {
					// Si no se espera un error específico, podría ser nil o io.EOF si se llegó al final.
					// io.EOF no es realmente un "error" en el uso normal de scanner.
					if err != nil && err != io.EOF {
						t.Errorf("Error inesperado del scanner: %v", err)
					}
				}
			}
		})
	}
}
