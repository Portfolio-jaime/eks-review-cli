# Instalación Avanzada y Acceso Global a `eks-review`

Esta guía explica cómo compilar tu herramienta CLI con el nombre `eks-review` y cómo instalarla para que puedas ejecutarla directamente desde cualquier ubicación en tu terminal, sin necesidad de navegar a su directorio de código fuente o usar `./eks-review`.

---

## 1. Cambiar el Nombre del Ejecutable Durante la Compilación

Por defecto, al compilar tu proyecto, el nombre del ejecutable podría basarse en el nombre del directorio o del módulo. Para asegurarte de que el binario se llame `eks-review`, usa la opción `-o`:

```bash
go build -o eks-review
```

Esto generará un archivo ejecutable llamado `eks-review` (o `eks-review.exe` en Windows) en el directorio actual del proyecto.

---

## 2. Instalar `eks-review` para Acceso Global

Para ejecutar `eks-review` desde cualquier lugar, el archivo ejecutable debe estar ubicado en un directorio que esté listado en la variable de entorno `PATH` de tu sistema.

### Opción A: Instalación Manual (Recomendada)

1. **Compila `eks-review`:**

    ```bash
    go build -o eks-review
    ```

2. **Mueve el ejecutable a un directorio en tu `PATH`:**

    - **Linux y macOS:**
        - Un directorio común es `/usr/local/bin` (requiere permisos de superusuario):

            ```bash
            sudo mv eks-review /usr/local/bin/
            ```

        - Alternativamente, puedes usar un directorio personal como `$HOME/bin` o `$HOME/.local/bin`. Asegúrate de que esté en tu `PATH`:

            ```bash
            mkdir -p $HOME/bin
            mv eks-review $HOME/bin/
            export PATH="$HOME/bin:$PATH"
            ```

            Añade la línea anterior a tu `~/.bashrc` o `~/.zshrc` para que sea permanente.

    - **Windows:**
        - Crea un directorio (ej. `C:\Program Files\eks-review` o `C:\Users\TuUsuario\bin`).
        - Mueve el archivo `eks-review.exe` a ese directorio.
        - Añade ese directorio a tu variable de entorno `PATH` a través de las "Variables de Entorno" en las propiedades del sistema.

3. **Verifica la instalación:**

    Abre una nueva terminal y ejecuta:

    ```bash
    eks-review --help
    ```

    Debería mostrar la ayuda de tu CLI.

---

### Opción B: Usando `go install`

El comando `go install` compila e instala paquetes. Los binarios se colocan en el directorio `$GOPATH/bin` (o `$HOME/go/bin` si `GOPATH` no está definido).

1. **Asegúrate de que el directorio de binarios de Go esté en tu `PATH`:**

    ```bash
    export PATH="$(go env GOPATH)/bin:$PATH"
    # O, si GOPATH no está definido explícitamente:
    # export PATH="$HOME/go/bin:$PATH"
    ```

    Añade la línea anterior a tu archivo de configuración de shell.

2. **Ejecuta `go install`:**

    Desde la raíz de tu proyecto:

    ```bash
    go install
    ```

    El nombre del binario será el del módulo o directorio principal. Si quieres que se llame `eks-review`, es más directo usar la Opción A.

---

### Opción C: Distribución a Otros Usuarios (Avanzado)

Si deseas distribuir `eks-review` a usuarios que no necesariamente tienen Go instalado:

1. **Compilación Cruzada:**

    ```bash
    # Para Linux AMD64
    GOOS=linux GOARCH=amd64 go build -o eks-review_linux_amd64
    # Para Windows AMD64
    GOOS=windows GOARCH=amd64 go build -o eks-review_windows_amd64.exe
    # Para macOS AMD64
    GOOS=darwin GOARCH=amd64 go build -o eks-review_darwin_amd64
    # Para macOS ARM64 (Apple Silicon)
    GOOS=darwin GOARCH=arm64 go build -o eks-review_darwin_arm64
    ```

2. **Empaquetado:**
    - Crea archivos `.zip` o `.tar.gz` que contengan el binario, `README.md`, `COMMANDS.md`, y otros archivos necesarios.
    - Considera incluir un script de instalación o instrucciones claras.

3. **Releases de GitHub:**
    - Sube estos archivos empaquetados y/o los binarios directamente a la sección "Releases" de tu repositorio de GitHub.

4. **Gestores de Paquetes (más avanzado):**
    - Homebrew (macOS/Linux), Scoop (Windows), paquetes `.deb` o `.rpm` para Linux, etc.

---

## Recomendación Inicial

Para desarrollo propio y compartir con otros desarrolladores Go, la Opción A (compilación manual y mover a un directorio en el `PATH`) o una variante de la Opción B son las más directas y rápidas.

Para distribución más amplia, empieza con Releases de GitHub (Opción C.3) proporcionando binarios precompilados.

---

## 📚 Integración sugerida en el README

Puedes añadir una sección como esta en tu `README.md` principal:

```markdown
## ⚙️ Instalación Avanzada (Acceso Global como `eks-review`)

Por defecto, después de compilar con `go build -o eks-review`, puedes ejecutar la herramienta desde el directorio del proyecto con `./eks-review`.

Si deseas poder ejecutar `eks-review` desde cualquier ubicación en tu terminal, necesitarás instalar el binario en un directorio que esté en tu `PATH` del sistema.

Para instrucciones detalladas sobre cómo compilar con el nombre `eks-review` e instalarlo globalmente en diferentes sistemas operativos, consulta nuestra [Guía de Instalación Avanzada](INSTALLATION_ADVANCED.md).
```
